package posts

import (
	"database/sql"

	"post-pilot/apps/api/internal/posts/handler"
	"post-pilot/apps/api/internal/posts/repository"
	"post-pilot/apps/api/internal/posts/service"
)

// Module wires the posts domain: repository → service → handler.
type Module struct {
	Handler *handler.Handler
	Service service.PostService
}

// NewModule wires up a posts Module.
func NewModule(db *sql.DB) *Module {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	return &Module{Handler: h, Service: svc}
}
