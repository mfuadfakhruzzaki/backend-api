// routes/userRoutes.go
package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/controllers"
	"github.com/mfuadfakhruzzaki/backend-api/middleware"
)

// RegisterUserRoutes registers all user-related routes under the provided router group
func RegisterUserRoutes(router *gin.RouterGroup) {
    fmt.Println("Registering user routes")

    // Create a subgroup for user-related routes
    userGroup := router.Group("/user")
    {
        // Apply JWT middleware to all routes in this group
        userGroup.Use(middleware.JWTMiddleware())

        // Endpoint untuk meng-upload profile picture
        userGroup.POST("/profile-picture", controllers.UploadProfilePicture)

        // Endpoint untuk mendapatkan profil pengguna
        userGroup.GET("/profile", controllers.GetProfile)
    }
}
