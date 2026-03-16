package bootstrap

import (
	"post-pilot/apps/api/internal/config"
	"post-pilot/packages/database"
)

type Container struct {
	DB     *database.DB
	Router *Router
	Logger Logger
	Config *config.Config
}

func NewContainer(logger Logger, cfg *config.Config) (*Container, error) {
	db, err := database.NewDB(cfg.DBName, cfg.DBUser, cfg.DBPassword)
	if err != nil {
		return nil, err
	}

	return &Container{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}, nil
}
