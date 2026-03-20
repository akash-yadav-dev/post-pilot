package social

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrSocialNotFound = errors.New("social record not found")

// type accountRow struct {
// 	SocialAccount
// 	AccessToken  string
// 	RefreshToken string
// 	Metadata     string
// }

type publishTarget struct {
	TargetID         uuid.UUID
	PostID           uuid.UUID
	PostUserID       uuid.UUID
	Content          string
	MediaURLs        []string
	Platform         string
	SocialAccountID  uuid.UUID
	ExternalAccount  string
	AccessToken      string
	RefreshToken     string
	Metadata         string
	AccessTokenValid bool
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertAccount(ctx context.Context, userID uuid.UUID, req ConnectAccountRequest) (*SocialAccount, error) {
	query := `
		INSERT INTO social_accounts (
			user_id, platform, account_id, account_name, account_url,
			access_token, refresh_token, token_expires_at, token_scope, metadata, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9, COALESCE(NULLIF($10, ''), '{}'::jsonb), 'active')
		ON CONFLICT (platform, account_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			account_name = EXCLUDED.account_name,
			account_url = EXCLUDED.account_url,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			token_expires_at = EXCLUDED.token_expires_at,
			token_scope = EXCLUDED.token_scope,
			metadata = EXCLUDED.metadata,
			status = 'active',
			last_error = NULL,
			error_count = 0,
			updated_at = NOW()
		RETURNING id, user_id, platform::text, account_id, COALESCE(account_name, ''), COALESCE(account_url, ''), token_expires_at, COALESCE(token_scope, ''), status::text, last_used_at, COALESCE(last_error, ''), error_count, created_at, updated_at
	`

	account := &SocialAccount{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		strings.ToLower(strings.TrimSpace(req.Platform)),
		strings.TrimSpace(req.AccountID),
		strings.TrimSpace(req.AccountName),
		strings.TrimSpace(req.AccountURL),
		strings.TrimSpace(req.AccessToken),
		strings.TrimSpace(req.RefreshToken),
		req.TokenExpiresAt,
		strings.TrimSpace(req.TokenScope),
		strings.TrimSpace(req.Metadata),
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Platform,
		&account.AccountID,
		&account.AccountName,
		&account.AccountURL,
		&account.TokenExpiresAt,
		&account.TokenScope,
		&account.Status,
		&account.LastUsedAt,
		&account.LastError,
		&account.ErrorCount,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *Repository) ListAccountsByUser(ctx context.Context, userID uuid.UUID) ([]*SocialAccount, error) {
	query := `
		SELECT id, user_id, platform::text, account_id, COALESCE(account_name, ''), COALESCE(account_url, ''), token_expires_at, COALESCE(token_scope, ''), status::text, last_used_at, COALESCE(last_error, ''), error_count, created_at, updated_at
		FROM social_accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*SocialAccount, 0)
	for rows.Next() {
		acc := &SocialAccount{}
		if err := rows.Scan(
			&acc.ID,
			&acc.UserID,
			&acc.Platform,
			&acc.AccountID,
			&acc.AccountName,
			&acc.AccountURL,
			&acc.TokenExpiresAt,
			&acc.TokenScope,
			&acc.Status,
			&acc.LastUsedAt,
			&acc.LastError,
			&acc.ErrorCount,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *Repository) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	query := `
		DELETE FROM social_accounts
		WHERE id = $1 AND user_id = $2
	`
	res, err := r.db.ExecContext(ctx, query, accountID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSocialNotFound
	}
	return nil
}

func (r *Repository) ListPublishTargets(ctx context.Context, userID, postID uuid.UUID) ([]*publishTarget, error) {
	query := `
		SELECT
			t.id,
			t.post_id,
			p.user_id,
			COALESCE(t.content_override, p.content),
			COALESCE(t.media_urls_override, p.media_urls),
			sa.platform::text,
			t.social_account_id,
			sa.account_id,
			sa.access_token,
			COALESCE(sa.refresh_token, ''),
			COALESCE(sa.metadata::text, '{}')
		FROM post_targets t
		INNER JOIN posts p ON p.id = t.post_id
		INNER JOIN social_accounts sa ON sa.id = t.social_account_id
		WHERE t.post_id = $1
			AND p.user_id = $2
			AND p.deleted_at IS NULL
			AND t.status IN ('pending', 'failed')
		ORDER BY t.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, postID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	targets := make([]*publishTarget, 0)
	for rows.Next() {
		t := &publishTarget{}
		if err := rows.Scan(
			&t.TargetID,
			&t.PostID,
			&t.PostUserID,
			&t.Content,
			&t.MediaURLs,
			&t.Platform,
			&t.SocialAccountID,
			&t.ExternalAccount,
			&t.AccessToken,
			&t.RefreshToken,
			&t.Metadata,
		); err != nil {
			return nil, err
		}

		t.AccessTokenValid = strings.TrimSpace(t.AccessToken) != ""
		targets = append(targets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		exists, err := r.postExistsForUser(ctx, userID, postID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrSocialNotFound
		}
	}

	return targets, nil
}

func (r *Repository) MarkTargetQueued(ctx context.Context, targetID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE post_targets SET status = 'publishing', attempts = attempts + 1, updated_at = NOW() WHERE id = $1`, targetID)
	return err
}

func (r *Repository) MarkTargetPublished(ctx context.Context, targetID uuid.UUID, externalID, externalURL string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE post_targets
		SET status = 'published', platform_post_id = $2, platform_post_url = $3, published_at = NOW(), last_error = NULL, updated_at = NOW()
		WHERE id = $1
	`, targetID, strings.TrimSpace(externalID), strings.TrimSpace(externalURL))
	return err
}

func (r *Repository) MarkTargetFailed(ctx context.Context, targetID uuid.UUID, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE post_targets
		SET status = 'failed', last_error = $2, next_attempt_at = NOW() + INTERVAL '10 minutes', updated_at = NOW()
		WHERE id = $1
	`, targetID, strings.TrimSpace(errMsg))
	return err
}

func (r *Repository) MarkSocialAccountSuccess(ctx context.Context, accountID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE social_accounts
		SET last_used_at = NOW(), status = 'active', last_error = NULL, error_count = 0, updated_at = NOW()
		WHERE id = $1
	`, accountID)
	return err
}

func (r *Repository) MarkSocialAccountFailure(ctx context.Context, accountID uuid.UUID, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE social_accounts
		SET status = 'error', last_error = $2, error_count = error_count + 1, updated_at = NOW()
		WHERE id = $1
	`, accountID, strings.TrimSpace(errMsg))
	return err
}

func (r *Repository) FinalizePostStatus(ctx context.Context, postID uuid.UUID) (string, error) {
	var failed int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM post_targets WHERE post_id = $1 AND status = 'failed'`, postID).Scan(&failed); err != nil {
		return "", err
	}

	status := "published"
	if failed > 0 {
		status = "failed"
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE posts
		SET status = $2,
			published_at = CASE WHEN $2 = 'published' THEN NOW() ELSE published_at END,
			updated_at = NOW()
		WHERE id = $1
	`, postID, status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *Repository) postExistsForUser(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM posts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) setPublishingStatus(ctx context.Context, postID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE posts SET status = 'publishing', updated_at = NOW() WHERE id = $1`, postID)
	return err
}

func (r *Repository) TouchTargetRetry(ctx context.Context, targetID uuid.UUID, retryAfter time.Duration) error {
	_, err := r.db.ExecContext(ctx, `UPDATE post_targets SET next_attempt_at = NOW() + $2::interval, updated_at = NOW() WHERE id = $1`, targetID, retryAfter.String())
	return err
}
