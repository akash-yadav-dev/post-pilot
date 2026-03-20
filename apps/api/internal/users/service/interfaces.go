package service

import (
	"context"

	"post-pilot/apps/api/internal/users/model"

	"github.com/google/uuid"
)

// UserRepository is the persistence contract expected by the service layer.
type UserRepository interface {
	Create(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserService is the business-logic contract consumed by the handler layer.
type UserService interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
