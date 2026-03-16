package bootstrap

import (
	"net/http"
	"post-pilot/apps/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	UserRouter    http.Handler
	ProjectRouter http.Handler
	// Add more routers as needed
}

func SetupRouter(container *Container) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Cors configuration
	cors := middleware.DefaultCORSConfig()
	router.Use(middleware.NewCORSMiddleware(cors))

	// health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}
