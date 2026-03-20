package auth

import (
	"database/sql"
	"time"

	"post-pilot/apps/api/internal/auth/handler"
	"post-pilot/apps/api/internal/auth/repository"
	"post-pilot/apps/api/internal/auth/service"
	"post-pilot/apps/api/internal/config"
	"post-pilot/apps/api/internal/middleware"
	rediscache "post-pilot/packages/cache/redis"
	"post-pilot/packages/logger"
	"post-pilot/packages/security"

	"github.com/gin-gonic/gin"
)

type Module struct {
	Handler           *handler.AuthHandler
	Service           *service.AuthService
	JWTService        *security.JWTService
	AuthRequired      gin.HandlerFunc
	LoginIdentity     gin.HandlerFunc
	LoginRateLimit    gin.HandlerFunc
	LoginEmailLimit   gin.HandlerFunc
	RegisterRateLimit gin.HandlerFunc
	RefreshRateLimit  gin.HandlerFunc
}

func NewModule(db *sql.DB, cfg *config.Config, cacheSvc *rediscache.CacheService, appLogger logger.Logger) (*Module, error) {
	tokenStore := NewRedisTokenStore(cacheSvc)
	authRepo := repository.NewAuthRepository(db)

	passwordSvc, err := security.NewPasswordService(cfg.PasswordHashCost)
	if err != nil {
		return nil, err
	}

	jwtSvc, err := security.NewJWTService(security.JWTConfig{
		AccessSecret:           cfg.JWTAccessSecretKey,
		RefreshSecret:          cfg.JWTRefreshSecretKey,
		AccessTokenExpiration:  time.Duration(cfg.JWTExpiry) * time.Second,
		RefreshTokenExpiration: time.Duration(cfg.JWTRefreshExpiry) * time.Second,
	}, tokenStore)
	if err != nil {
		return nil, err
	}

	authService := service.NewAuthService(
		authRepo,
		passwordSvc,
		jwtSvc,
		service.NewGoogleIDTokenVerifier(cfg.GoogleOAuthClientID),
		time.Duration(cfg.JWTExpiry)*time.Second,
	)

	auditLogger := NewDBAuditLogger(db, appLogger)
	authHandler := handler.NewAuthHandler(authService, auditLogger)

	window := time.Duration(cfg.AuthRateLimitWindowSeconds) * time.Second

	return &Module{
		Handler:    authHandler,
		Service:    authService,
		JWTService: jwtSvc,

		AuthRequired:  middleware.NewAuthMiddleware(jwtSvc),
		LoginIdentity: middleware.NewLoginIdentityMiddleware(),
		LoginRateLimit: middleware.NewRateLimitMiddleware(cacheSvc, middleware.RateLimitConfig{
			Prefix: "auth:login:ip",
			Limit:  int64(cfg.AuthLoginRateLimitMax),
			Window: window,
			IdentifierResolver: func(c *gin.Context) string {
				return c.ClientIP()
			},
		}),
		LoginEmailLimit: middleware.NewRateLimitMiddleware(cacheSvc, middleware.RateLimitConfig{
			Prefix: "auth:login:email",
			Limit:  int64(cfg.AuthLoginRateLimitMax),
			Window: window,
			IdentifierResolver: func(c *gin.Context) string {
				email := c.GetString("auth_login_email")
				if email == "" {
					return "unknown"
				}
				return email
			},
		}),
		RegisterRateLimit: middleware.NewRateLimitMiddleware(cacheSvc, middleware.RateLimitConfig{
			Prefix: "auth:register:ip",
			Limit:  int64(cfg.AuthRegisterRateLimitMax),
			Window: window,
			IdentifierResolver: func(c *gin.Context) string {
				return c.ClientIP()
			},
		}),
		RefreshRateLimit: middleware.NewRateLimitMiddleware(cacheSvc, middleware.RateLimitConfig{
			Prefix: "auth:refresh:ip",
			Limit:  int64(cfg.AuthRefreshRateLimitMax),
			Window: window,
			IdentifierResolver: func(c *gin.Context) string {
				return c.ClientIP()
			},
		}),
	}, nil
}
