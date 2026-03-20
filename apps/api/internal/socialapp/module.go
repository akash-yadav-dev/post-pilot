package socialapp

import (
	"database/sql"

	"post-pilot/apps/api/internal/social"
	"post-pilot/apps/api/internal/social/bluesky"
	"post-pilot/apps/api/internal/social/linkedin"
	"post-pilot/apps/api/internal/social/mastodon"
	"post-pilot/apps/api/internal/social/twitter"
)

type Module struct {
	Handler *social.Handler
	Service *social.Service
}

func NewModule(db *sql.DB) *Module {
	repo := social.NewRepository(db)
	registry := social.NewRegistry()

	registry.Register(linkedin.NewPublisher(linkedin.NewClient("", "")))
	registry.Register(twitter.NewPublisher(twitter.NewClient("", "", "", "", "")))
	registry.Register(mastodon.NewPublisher(mastodon.NewClient("", "")))
	registry.Register(bluesky.NewPublisher(bluesky.NewClient("", "", "")))

	svc := social.NewService(repo, registry)
	h := social.NewHandler(svc)

	return &Module{Handler: h, Service: svc}
}
