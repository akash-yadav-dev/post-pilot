package posts

import (
	"context"
	"database/sql"

	"post-pilot/apps/api/internal/posts/handler"
	"post-pilot/apps/api/internal/posts/model"
	"post-pilot/apps/api/internal/posts/service"

	"github.com/google/uuid"
)

// Module wires the posts domain: repository → service → handler.
type Module struct {
	Handler *handler.Handler
	Service service.PostService
}

// NewModule wires up a posts Module.
// stubPostRepository is used until the full SQL repository is implemented.
func NewModule(db *sql.DB) *Module {
	var repo service.PostRepository = &stubPostRepository{db: db}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	return &Module{Handler: h, Service: svc}
}

// stubPostRepository satisfies service.PostRepository end-to-end at compile time.
// Every method panics so any accidental call is surfaced immediately during dev.
type stubPostRepository struct{ db *sql.DB }

func (r *stubPostRepository) Create(_ context.Context, _ uuid.UUID, _ model.CreatePostRequest) (*model.Post, error) {
	panic("posts: repository.Create not implemented")
}
func (r *stubPostRepository) GetByID(_ context.Context, _ uuid.UUID) (*model.Post, error) {
	panic("posts: repository.GetByID not implemented")
}
func (r *stubPostRepository) ListByUser(_ context.Context, _ uuid.UUID) ([]*model.Post, error) {
	panic("posts: repository.ListByUser not implemented")
}
func (r *stubPostRepository) Update(_ context.Context, _ uuid.UUID, _ model.UpdatePostRequest) (*model.Post, error) {
	panic("posts: repository.Update not implemented")
}
func (r *stubPostRepository) UpdateStatus(_ context.Context, _ uuid.UUID, _ model.PostStatus) error {
	panic("posts: repository.UpdateStatus not implemented")
}
func (r *stubPostRepository) Delete(_ context.Context, _ uuid.UUID) error {
	panic("posts: repository.Delete not implemented")
}
