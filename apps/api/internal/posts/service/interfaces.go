package service

import (
	"context"

	"post-pilot/apps/api/internal/posts/model"

	"github.com/google/uuid"
)

// PostRepository is the persistence contract expected by the service layer.
type PostRepository interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PostStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostService is the business-logic contract consumed by the handler layer.
type PostService interface {
	CreatePost(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	ListUserPosts(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	UpdatePost(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error)
	SchedulePost(ctx context.Context, id uuid.UUID) error
	DeletePost(ctx context.Context, id uuid.UUID) error
}
