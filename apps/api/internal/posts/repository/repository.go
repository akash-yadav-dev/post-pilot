package repository

import (
	"context"

	"post-pilot/apps/api/internal/posts/model"

	"github.com/google/uuid"
)

// PostRepository defines persistence operations for the posts domain.
type PostRepository interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CreatePostRequest) (*model.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Post, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdatePostRequest) (*model.Post, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PostStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
