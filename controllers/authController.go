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
// @Summary Register a new user
// @Description This endpoint allows users to register by providing email, username, password, and phone number. A verification email will be sent after registration.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user  body     models.User  true  "User data"
// @Success 201 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]interface{} "Invalid request payload or email/username already exists"
// @Failure 500 {object} map[string]interface{} "Error creating user or sending verification email"
// @Router  /auth/register [post]
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

	// Generate verification code
	verificationCode := utils.GenerateVerificationCode()

	// Create new user with hashed password and verification code
	user := models.User{
		Email:            userInput.Email,
		Username:         userInput.Username,
		Password:         hashedPassword,
		PhoneNumber:      userInput.PhoneNumber,
		ProfilePicture:   "", // Initialize with empty string
		PackageID:        nil,
		EmailVerified:    false, // Email not verified yet
		VerificationCode: verificationCode,
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

	// Send verification email
	if err := utils.SendVerificationEmail(user.Email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// Remove password before sending response
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "Registration successful! Please check your email to verify your account.",
	})
}

// VerifyEmail handles the verification of user's email
// @Summary Verify user email
// @Description This endpoint allows users to verify their email by providing the verification code sent via email.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   verification  body  map[string]string  true  "Email and verification code"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload or verification code"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Failed to verify email"
// @Router  /auth/verify [post]
func VerifyEmail(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	// Binding JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var user models.User
	// Find user by email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the verification code is correct
	if user.VerificationCode != input.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Update user to set email as verified
	user.EmailVerified = true
	user.VerificationCode = "" // Optionally clear the verification code
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
}

// Login handles user authentication
// @Summary User login
// @Description This endpoint allows users to log in by providing email and password. A JWT token will be returned upon successful login.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials  body  map[string]string  true  "User credentials (email and password)"
// @Success 200 {object} map[string]interface{} "JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized, invalid credentials or email not verified"
// @Failure 500 {object} map[string]interface{} "Error generating token or database error"
// @Router  /auth/login [post]
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

	// Check if email is verified
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified. Please verify your email first."})
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
