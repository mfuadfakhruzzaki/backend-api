// controllers/authController.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"github.com/mfuadfakhruzzaki/backend-api/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Register handles user registration
func Register(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure password is not empty
	if strings.TrimSpace(userInput.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create new user with hashed password
	user := models.User{
		Email:          userInput.Email,
		Username:       userInput.Username,
		Password:       hashedPassword,
		PhoneNumber:    userInput.PhoneNumber,
		ProfilePicture: "", // Initialize with empty string
		PackageID:      nil,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		// Check for duplicate entry error (unique constraint violation)
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Remove password before sending response
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// Login handles user authentication
func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", credentials.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
