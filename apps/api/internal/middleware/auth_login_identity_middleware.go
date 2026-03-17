package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

type loginIdentityPayload struct {
	Email string `json:"email"`
}

func NewLoginIdentityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request == nil || c.Request.Body == nil {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(nil))
			c.Next()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var payload loginIdentityPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			c.Next()
			return
		}

		normalizedEmail := strings.ToLower(strings.TrimSpace(payload.Email))
		if normalizedEmail != "" {
			c.Set("auth_login_email", normalizedEmail)
		}

		c.Next()
	}
}
