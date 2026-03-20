package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"post-pilot/apps/api/internal/users/model"
	"post-pilot/apps/api/internal/users/repository"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already registered")
)

// Ensure Service satisfies UserService at compile time.
var _ UserService = (*Service)(nil)

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.Update(ctx, id, req)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrUserNotFound
	} else if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
