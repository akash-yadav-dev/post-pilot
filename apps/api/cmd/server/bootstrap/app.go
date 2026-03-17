package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"post-pilot/apps/api/internal/config"
	"syscall"
	"time"
)

type App struct {
	Container *Container
	Router    http.Handler
}

func NewApp() (*App, error) {
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

	router, err := SetupRouter(container)
	if err != nil {
		return nil, fmt.Errorf("failed to setup router: %w", err)
	}

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

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           a.Router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("Server running on port %s", port)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
