package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"post-pilot/apps/api/internal/auth/model"

	"github.com/google/uuid"
)

var (
	ErrNotFound           = errors.New("record not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUserWithPassword(ctx context.Context, name, email, passwordHash string) (*model.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	user := &model.User{}

	insertUserQuery := `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id, email, name
	`

	err = tx.QueryRowContext(ctx, insertUserQuery, email, name).
		Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	insertAuthQuery := `
		INSERT INTO auth_accounts (
			user_id,
			provider,
			provider_user_id,
			password_hash,
			password_changed_at
		)
		VALUES ($1, 'password', $2, $3, NOW())
	`

	_, err = tx.ExecContext(ctx, insertAuthQuery, user.ID, email, passwordHash)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) FindUserByProviderIdentity(ctx context.Context, provider, providerUserID string) (*model.User, error) {
	query := `
		SELECT u.id, u.email, u.name
		FROM users u
		INNER JOIN auth_accounts a ON a.user_id = u.id
		WHERE a.provider = $1
			AND a.provider_user_id = $2
			AND u.deleted_at IS NULL
			AND u.status = 'active'
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) CreateOrLinkGoogleUser(ctx context.Context, name, email, providerUserID string, emailVerified bool) (*model.User, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	user := &model.User{}
	findByEmailQuery := `
		SELECT id, email, name
		FROM users
		WHERE email = $1
			AND deleted_at IS NULL
		LIMIT 1
	`

	err = tx.QueryRowContext(ctx, findByEmailQuery, email).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		createUserQuery := `
			INSERT INTO users (email, name, email_verified_at)
			VALUES ($1, $2, CASE WHEN $3 THEN NOW() ELSE NULL END)
			RETURNING id, email, name
		`
		if err := tx.QueryRowContext(ctx, createUserQuery, email, name, emailVerified).Scan(&user.ID, &user.Email, &user.Name); err != nil {
			if isUniqueViolation(err) {
				return nil, ErrEmailAlreadyExists
			}
			return nil, err
		}
	}

	upsertAuthQuery := `
		INSERT INTO auth_accounts (user_id, provider, provider_user_id)
		VALUES ($1, 'google', $2)
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET user_id = EXCLUDED.user_id, updated_at = NOW()
	`

	if _, err := tx.ExecContext(ctx, upsertAuthQuery, user.ID, providerUserID); err != nil {
		return nil, err
	}

	if emailVerified {
		markVerifiedQuery := `
			UPDATE users
			SET email_verified_at = COALESCE(email_verified_at, NOW())
			WHERE id = $1
		`
		if _, err := tx.ExecContext(ctx, markVerifiedQuery, user.ID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) GetPasswordIdentityByEmail(ctx context.Context, email string) (*model.PasswordAuthIdentity, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.name,
			a.password_hash,
			a.failed_login_count,
			a.locked_until
		FROM users u
		INNER JOIN auth_accounts a ON a.user_id = u.id
		WHERE
			a.provider = 'password'
			AND a.provider_user_id = $1
			AND u.deleted_at IS NULL
			AND u.status = 'active'
	`

	identity := &model.PasswordAuthIdentity{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&identity.ID,
		&identity.Email,
		&identity.Name,
		&identity.PasswordHash,
		&identity.FailedCount,
		&identity.LockedUntil,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return identity, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, name
		FROM users
		WHERE id = $1
		AND deleted_at IS NULL
		AND status = 'active'
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) RecordFailedLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE auth_accounts
		SET
			failed_login_count = failed_login_count + 1,
			locked_until = CASE
				WHEN failed_login_count + 1 >= 5 THEN NOW() + INTERVAL '15 minutes'
				ELSE locked_until
			END,
			updated_at = NOW()
		WHERE user_id = $1
		AND provider = 'password'
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *AuthRepository) ResetFailedLogin(ctx context.Context, userID uuid.UUID) error {

	query := `
		UPDATE auth_accounts
		SET
			failed_login_count = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE user_id = $1
		AND provider = 'password'
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *AuthRepository) UpdatePasswordHash(ctx context.Context, userID uuid.UUID, newHash string) error {
	query := `
		UPDATE auth_accounts
		SET
			password_hash = $2,
			password_changed_at = NOW(),
			updated_at = NOW()
		WHERE user_id = $1
		AND provider = 'password'
	`

	_, err := r.db.ExecContext(ctx, query, userID, newHash)
	return err
}

func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "unique constraint")
}
