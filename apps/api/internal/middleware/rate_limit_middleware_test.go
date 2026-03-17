package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rediscache "post-pilot/packages/cache/redis"

	"github.com/gin-gonic/gin"
)

type fakeRateLimitRepo struct {
	counters map[string]int64
}

func newFakeRateLimitRepo() *fakeRateLimitRepo {
	return &fakeRateLimitRepo{counters: map[string]int64{}}
}

func (f *fakeRateLimitRepo) Get(_ context.Context, _ string) (string, error) {
	return "", rediscache.ErrCacheMiss
}

func (f *fakeRateLimitRepo) Set(_ context.Context, _ string, _ string, _ time.Duration) error {
	return nil
}

func (f *fakeRateLimitRepo) Delete(_ context.Context, key string) error {
	delete(f.counters, key)
	return nil
}

func (f *fakeRateLimitRepo) Exists(_ context.Context, key string) (bool, error) {
	_, ok := f.counters[key]
	return ok, nil
}

func (f *fakeRateLimitRepo) Increment(_ context.Context, key string, _ time.Duration) (int64, error) {
	f.counters[key]++
	return f.counters[key], nil
}

func (f *fakeRateLimitRepo) HealthCheck(_ context.Context) error {
	return nil
}

type failingRateLimitRepo struct{}

func (f *failingRateLimitRepo) Get(_ context.Context, _ string) (string, error) {
	return "", errors.New("redis down")
}
func (f *failingRateLimitRepo) Set(_ context.Context, _ string, _ string, _ time.Duration) error {
	return errors.New("redis down")
}
func (f *failingRateLimitRepo) Delete(_ context.Context, _ string) error {
	return errors.New("redis down")
}
func (f *failingRateLimitRepo) Exists(_ context.Context, _ string) (bool, error) {
	return false, errors.New("redis down")
}
func (f *failingRateLimitRepo) Increment(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, errors.New("redis down")
}
func (f *failingRateLimitRepo) HealthCheck(_ context.Context) error {
	return errors.New("redis down")
}

func TestRateLimitMiddlewareBlocksAfterLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cacheSvc := rediscache.NewCacheService(newFakeRateLimitRepo(), rediscache.CacheOptions{DefaultTTL: time.Minute})
	limiter := NewRateLimitMiddleware(cacheSvc, RateLimitConfig{
		Prefix: "auth:login:ip",
		Limit:  2,
		Window: time.Minute,
		IdentifierResolver: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})

	router := gin.New()
	router.POST("/login", limiter, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, w.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimitMiddlewareCombinedByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cacheSvc := rediscache.NewCacheService(newFakeRateLimitRepo(), rediscache.CacheOptions{DefaultTTL: time.Minute})
	emailLimiter := NewRateLimitMiddleware(cacheSvc, RateLimitConfig{
		Prefix: "auth:login:email",
		Limit:  1,
		Window: time.Minute,
		IdentifierResolver: func(c *gin.Context) string {
			return c.GetString("auth_login_email")
		},
	})

	router := gin.New()
	router.POST("/login", func(c *gin.Context) {
		c.Set("auth_login_email", c.Query("email"))
		c.Next()
	}, emailLimiter, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/login?email=a@example.com", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/login?email=a@example.com", nil))
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}

	third := httptest.NewRecorder()
	router.ServeHTTP(third, httptest.NewRequest(http.MethodPost, "/login?email=b@example.com", nil))
	if third.Code != http.StatusOK {
		t.Fatalf("third status = %d, want %d", third.Code, http.StatusOK)
	}
}

func TestRateLimitMiddlewareHandlesBackendFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cacheSvc := rediscache.NewCacheService(&failingRateLimitRepo{}, rediscache.CacheOptions{DefaultTTL: time.Minute})
	limiter := NewRateLimitMiddleware(cacheSvc, RateLimitConfig{
		Prefix: "auth:login:ip",
		Limit:  1,
		Window: time.Minute,
		IdentifierResolver: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})

	router := gin.New()
	router.POST("/login", limiter, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
