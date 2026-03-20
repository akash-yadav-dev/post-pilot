package service

import (
	"context"
	"errors"
	"fmt"

	"post-pilot/apps/api/internal/posts/model"

	"github.com/google/uuid"
)

var ErrPostNotFound = errors.New("post not found")

// Ensure Service satisfies PostService at compile time.
var _ PostService = (*Service)(nil)

type Service struct {
	repo PostRepository
}

func NewService(repo PostRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error) {
	post, err := s.repo.Create(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return post, nil
}

func (s *Service) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	return post, nil
}

func (s *Service) ListUserPosts(ctx context.Context, userID uuid.UUID) ([]*model.Post, error) {
	posts, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	return posts, nil
}

func (s *Service) UpdatePost(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error) {
	post, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	return post, nil
}

func (s *Service) SchedulePost(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UpdateStatus(ctx, id, model.PostStatusScheduled); err != nil {
		return fmt.Errorf("schedule post: %w", err)
	}
	return nil
}

func (s *Service) DeletePost(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}
