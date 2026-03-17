package service

import (
	"context"
	"post-pilot/apps/api/internal/auth/model"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUserWithPassword(ctx context.Context, name, email, passwordHash string) (*model.User, error)
	GetPasswordIdentityByEmail(ctx context.Context, email string) (*model.PasswordAuthIdentity, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	RecordFailedLogin(ctx context.Context, userID uuid.UUID) error
	ResetFailedLogin(ctx context.Context, userID uuid.UUID) error
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, newHash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}
