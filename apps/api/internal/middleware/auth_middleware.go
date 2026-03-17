package middleware

import (
	"errors"
	"net/http"
	"strings"

	"post-pilot/packages/security"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			message := "unauthorized"

			switch {
			case errors.Is(err, security.ErrExpiredToken):
				message = "access token has expired"
			case errors.Is(err, security.ErrWrongTokenType):
				message = "invalid token type"
			}

			c.AbortWithStatusJSON(status, gin.H{"error": message})
			return
		}

		c.Set("auth_user_id", claims.UserID)
		c.Set("auth_email", claims.Email)
		c.Next()
	}
}
