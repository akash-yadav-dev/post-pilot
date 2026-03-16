package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for CORS middleware
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
}

// DefaultCORSConfig returns production ready defaults
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000", // dev
			"http://localhost:3001",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
	}
}

// NewCORSMiddleware creates a CORS middleware
func NewCORSMiddleware(config CORSConfig) gin.HandlerFunc {

	// Allow overriding via ENV
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		config.AllowedOrigins = strings.Split(envOrigins, ",")
	}

	allowedOrigins := make(map[string]bool)
	for _, origin := range config.AllowedOrigins {
		allowedOrigins[strings.TrimSpace(origin)] = true
	}

	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		if origin != "" && allowedOrigins[origin] {

			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set(
				"Access-Control-Allow-Methods",
				strings.Join(config.AllowedMethods, ","),
			)

			c.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				strings.Join(config.AllowedHeaders, ","),
			)

			c.Writer.Header().Set(
				"Access-Control-Expose-Headers",
				strings.Join(config.ExposeHeaders, ","),
			)

			if config.AllowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		// Handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
