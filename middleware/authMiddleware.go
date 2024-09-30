package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/utils"
)

type ContextKey string

const (
	UserContextKey ContextKey = "userEmail"
	AuthHeader     string     = "Authorization"
	BearerSchema   string     = "bearer"
)

// JWTMiddleware verifies the JWT token and adds the user's email to the Gin context
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authHeader := c.GetHeader(AuthHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Split the header to extract the token
		// Expected format: "Bearer <token>"
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != BearerSchema {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		// Extract the token part
		tokenString := tokenParts[1]

		// Validate the token and extract the email
		email, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		// Store the email in the Gin context
		c.Set(string(UserContextKey), email)

		// Proceed to the next middleware or handler
		c.Next()
	}
}
