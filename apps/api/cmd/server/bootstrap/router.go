package bootstrap

import (
	"fmt"
	"net/http"
	"post-pilot/apps/api/internal/auth"
	"post-pilot/apps/api/internal/middleware"
	"post-pilot/apps/api/internal/posts"
	"post-pilot/apps/api/internal/users"
	"post-pilot/apps/api/routes"
	rediscache "post-pilot/packages/cache/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *Container) (*gin.Engine, error) {
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

	redisClient, err := rediscache.NewRedisClient()
	if err != nil {
		return nil, fmt.Errorf("setup redis client: %w", err)
	}

	redisRepository := rediscache.NewRepository(redisClient)
	cacheService := rediscache.NewCacheService(redisRepository, rediscache.CacheOptions{
		DefaultTTL: 5 * time.Minute,
		KeyPrefix:  "post_pilot",
	})

	authModule, err := auth.NewModule(container.DB.DB, container.Config, cacheService, container.Logger)
	if err != nil {
		return nil, fmt.Errorf("setup auth module: %w", err)
	}

	usersModule := users.NewModule(container.DB.DB)
	postsModule := posts.NewModule(container.DB.DB)

	routes.SetupAuthRouter(router, authModule)
	routes.SetupUserRouter(router, usersModule, authModule.AuthRequired)
	routes.SetupPostRouter(router, postsModule, authModule.AuthRequired)

	return router, nil
}
