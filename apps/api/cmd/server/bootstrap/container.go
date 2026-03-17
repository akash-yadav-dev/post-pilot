package bootstrap

import (
	"post-pilot/apps/api/internal/config"
	"post-pilot/packages/database/postgres"
)

type Container struct {
	DB     *postgres.DB
	Logger Logger
	Config *config.Config
}

func NewContainer(logger Logger, cfg *config.Config) (*Container, error) {
	db, err := postgres.NewPostgresDB()
	if err != nil {
		return nil, err
	}

	return &Container{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}, nil
}
