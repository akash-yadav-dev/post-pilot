package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginIdentityMiddlewareSetsNormalizedEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/login", NewLoginIdentityMiddleware(), func(c *gin.Context) {
		var payload map[string]any
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"email":   c.GetString("auth_login_email"),
			"payload": payload,
		})
	})

	body := `{"email":"  USER@Example.COM ","password":"secret"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Email string         `json:"email"`
		Data  map[string]any `json:"payload"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Email != "user@example.com" {
		t.Fatalf("email = %q, want %q", resp.Email, "user@example.com")
	}

	if _, ok := resp.Data["password"]; !ok {
		t.Fatal("expected payload to still contain password field after middleware")
	}
}

func TestLoginIdentityMiddlewareSkipsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/login", NewLoginIdentityMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"email": c.GetString("auth_login_email")})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if strings.Contains(w.Body.String(), "@") {
		t.Fatalf("unexpected email extraction for invalid json: %s", w.Body.String())
	}
}
