// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mfuadfakhruzzaki/backend-api/controllers"
	"github.com/mfuadfakhruzzaki/backend-api/middleware"
)

func RegisterRoutes(router *gin.Engine) {
	// Public Routes
	public := router.Group("/")
	{
		// Registration and Login Endpoints
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		// OAuth Google Endpoints
		public.GET("/auth/google/login", controllers.GoogleLogin)
		public.GET("/auth/google/callback", controllers.GoogleCallback)

		// OAuth GitHub Endpoints
		public.GET("/auth/github/login", controllers.GithubLogin)
		public.GET("/auth/github/callback", controllers.GithubCallback)
	}

	// Protected Routes with JWT Middleware
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware())
	{
		// Package Endpoints
		api.GET("/packages", controllers.GetPackages)
		api.POST("/packages/:id/select", controllers.SelectPackage)

		// User Endpoints
		api.POST("/user/profile-picture", controllers.UploadProfilePicture)
		api.GET("/user/profile", controllers.GetProfile)
	}
}
