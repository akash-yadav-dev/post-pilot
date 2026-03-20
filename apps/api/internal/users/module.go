package users

import (
	"database/sql"

	"post-pilot/apps/api/internal/users/handler"
	"post-pilot/apps/api/internal/users/repository"
	"post-pilot/apps/api/internal/users/service"
)

// Module wires the users domain: repository → service → handler.
type Module struct {
	Handler *handler.Handler
	Service service.UserService
}

func NewModule(db *sql.DB) *Module {
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	return &Module{
		Handler: h,
		Service: svc,
	}
}
