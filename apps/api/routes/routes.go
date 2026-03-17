package routes

import (
	"post-pilot/apps/api/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(router *gin.Engine, module *auth.Module) *gin.Engine {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", module.RegisterRateLimit, module.Handler.Register)
		authGroup.POST("/login", module.LoginIdentity, module.LoginRateLimit, module.LoginEmailLimit, module.Handler.Login)
		authGroup.POST("/refresh", module.RefreshRateLimit, module.Handler.Refresh)
		authGroup.POST("/logout", module.Handler.Logout)
	}

	protected := authGroup.Group("")
	protected.Use(module.AuthRequired)
	{
		protected.GET("/me", module.Handler.Me)
	}

	return router
}
