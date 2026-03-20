package repository

import (
	"context"
	"database/sql"
	"errors"

	"post-pilot/apps/api/internal/users/model"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("user not found")

// UserRepository defines persistence operations for the users domain.
type UserRepository interface {
	Create(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresRepository is the Postgres-backed implementation of UserRepository.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	const query = `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id, email, name, plan, credits, is_active, created_at, updated_at
	`
	u := &model.User{}
	err := r.db.QueryRowContext(ctx, query, req.Email, req.Name).
		Scan(&u.ID, &u.Email, &u.Name, &u.Plan, &u.Credits, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	const query = `
		SELECT id, email, name, plan, credits, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	u := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Email, &u.Name, &u.Plan, &u.Credits, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, email, name, plan, credits, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	u := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Name, &u.Plan, &u.Credits, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error) {
	const query = `
		UPDATE users
		SET name        = COALESCE(NULLIF($1, ''), name),
		    plan        = COALESCE(NULLIF($2, ''), plan),
		    updated_at  = NOW()
		WHERE id = $3
		RETURNING id, email, name, plan, credits, is_active, created_at, updated_at
	`
	u := &model.User{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Plan, id).
		Scan(&u.ID, &u.Email, &u.Name, &u.Plan, &u.Credits, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
