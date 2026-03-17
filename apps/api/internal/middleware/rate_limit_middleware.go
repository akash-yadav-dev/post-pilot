package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	rediscache "post-pilot/packages/cache/redis"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	Prefix             string
	Limit              int64
	Window             time.Duration
	IdentifierResolver func(c *gin.Context) string
}

func NewRateLimitMiddleware(cacheSvc *rediscache.CacheService, cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := strings.TrimSpace(cfg.IdentifierResolver(c))
		if identifier == "" {
			identifier = "unknown"
		}

		key := fmt.Sprintf("%s:%s", strings.TrimSpace(cfg.Prefix), identifier)

		count, err := cacheSvc.Increment(c.Request.Context(), key, cfg.Window)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiter unavailable"})
			return
		}

		if count > cfg.Limit {
			retryAfter := int(cfg.Window.Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
