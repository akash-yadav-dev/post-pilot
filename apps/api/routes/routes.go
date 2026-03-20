package routes

import (
	"post-pilot/apps/api/internal/auth"
	"post-pilot/apps/api/internal/posts"
	"post-pilot/apps/api/internal/users"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(router *gin.Engine, module *auth.Module) *gin.Engine {
	authGroup := router.Group("/api/v1/auth")
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

func SetupUserRouter(router *gin.Engine, module *users.Module, authMiddleware gin.HandlerFunc) *gin.Engine {
	userGroup := router.Group("/api/v1/users")
	userGroup.Use(authMiddleware)
	{
		userGroup.POST("", module.Handler.CreateUser)
		userGroup.GET("/:id", module.Handler.GetUser)
		userGroup.PATCH("/:id", module.Handler.UpdateUser)
		userGroup.DELETE("/:id", module.Handler.DeleteUser)
	}
	return router
}

func SetupPostRouter(router *gin.Engine, module *posts.Module, authMiddleware gin.HandlerFunc) *gin.Engine {
	postGroup := router.Group("/api/v1/posts")
	postGroup.Use(authMiddleware)
	{
		postGroup.POST("", module.Handler.CreatePost)
		postGroup.GET("", module.Handler.ListPosts)
		postGroup.GET("/:id", module.Handler.GetPost)
		postGroup.PATCH("/:id", module.Handler.UpdatePost)
		postGroup.DELETE("/:id", module.Handler.DeletePost)
	}
	return router
}
