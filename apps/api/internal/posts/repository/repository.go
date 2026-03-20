package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"post-pilot/apps/api/internal/posts/model"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("post not found")

// PostRepository defines persistence operations for the posts domain.
type PostRepository interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PostStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	status := model.PostStatusDraft
	if req.ScheduledAt != nil {
		status = model.PostStatusScheduled
	}

	post := &model.Post{}
	const insertPost = `
		INSERT INTO posts (user_id, content, status, scheduled_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, content, status, scheduled_at, published_at, created_at, updated_at
	`

	if err := tx.QueryRowContext(ctx, insertPost, userID, strings.TrimSpace(req.Content), status, req.ScheduledAt).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if err := r.insertTargetsForPlatforms(ctx, tx, post.ID, userID, req.Platforms); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	const query = `
		SELECT id, user_id, content, status, scheduled_at, published_at, created_at, updated_at
		FROM posts
		WHERE id = $1 AND deleted_at IS NULL
	`

	post := &model.Post{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Post, error) {
	const query = `
		SELECT id, user_id, content, status, scheduled_at, published_at, created_at, updated_at
		FROM posts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*model.Post, 0)
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.Status,
			&post.ScheduledAt,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error) {
	status := model.PostStatusDraft
	if req.ScheduledAt != nil {
		status = model.PostStatusScheduled
	}

	const query = `
		UPDATE posts
		SET
			content = COALESCE(NULLIF($1, ''), content),
			scheduled_at = $2,
			status = $3
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, user_id, content, status, scheduled_at, published_at, created_at, updated_at
	`

	post := &model.Post{}
	err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(req.Content), req.ScheduledAt, status, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Status,
		&post.ScheduledAt,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PostStatus) error {
	const query = `
		UPDATE posts
		SET status = $2,
			published_at = CASE WHEN $2 = 'published' THEN NOW() ELSE published_at END,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE posts
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) insertTargetsForPlatforms(ctx context.Context, tx *sql.Tx, postID, userID uuid.UUID, platforms []string) error {
	const selectAccount = `
		SELECT id
		FROM social_accounts
		WHERE user_id = $1
			AND platform = $2
			AND status = 'active'
		ORDER BY updated_at DESC
		LIMIT 1
	`

	const insertTarget = `
		INSERT INTO post_targets (post_id, social_account_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (post_id, social_account_id) DO NOTHING
	`

	for _, rawPlatform := range platforms {
		platform := strings.ToLower(strings.TrimSpace(rawPlatform))
		if platform == "" {
			continue
		}

		var socialAccountID uuid.UUID
		err := tx.QueryRowContext(ctx, selectAccount, userID, platform).Scan(&socialAccountID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no active social account connected for platform %q", platform)
		}
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, insertTarget, postID, socialAccountID); err != nil {
			return err
		}
	}

	return nil
}
