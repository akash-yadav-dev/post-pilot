package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"post-pilot/apps/api/internal/config"
)

type App struct {
	Container *Container
	Router    http.Handler
}

func NewApp() (*App, error) {
	// Initialize dependencies like database connection, router, etc. here
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	logger, err := NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	container, err := NewContainer(logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	router := SetupRouter(container)

	return &App{
		Container: container,
		Router:    router,
	}, nil
}

func (a *App) Start() {
	port := "8080"
	if a.Container != nil && a.Container.Config != nil && a.Container.Config.ServePort != "" {
		port = a.Container.Config.ServePort
	}

	log.Printf("Server is running on port %s", port)
	// Start the server at the specified port

	err := http.ListenAndServe(":"+port, a.Router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
