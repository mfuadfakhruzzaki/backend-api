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

// RegisterRequest represents the structure of the registration request body
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

// VerificationRequest represents the structure of the email verification request body
type VerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// LoginCredentials represents the structure of the login request body
type LoginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register handles user registration
// @Summary Register a new user
// @Description This endpoint allows users to register by providing email, username, password, and phone number. A verification email will be sent after registration.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user  body  RegisterRequest  true  "User registration data"
// @Success 201 {object} SuccessResponse "Registration successful"
// @Failure 400 {object} ErrorResponse "Invalid request payload or password is empty"
// @Failure 409 {object} ErrorResponse "Email or username already exists"
// @Failure 500 {object} ErrorResponse "Error creating user or sending verification email"
// @Router  /auth/register [post]
func Register(c *gin.Context) {
	var userInput RegisterRequest
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Ensure password is not empty
	if strings.TrimSpace(userInput.Password) == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Password cannot be empty"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error hashing password"})
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
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Email or username already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	// Send verification email
	if err := utils.SendVerificationEmail(user.Email, verificationCode); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to send verification email"})
		return
	}

	// Remove password before sending response
	user.Password = ""

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Registration successful! Please check your email to verify your account.",
		Data:    user,
	})
}

// VerifyEmail handles the verification of user's email
// @Summary Verify user email
// @Description This endpoint allows users to verify their email by providing the verification code sent via email.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   verification  body  VerificationRequest  true  "Email and verification code"
// @Success 200 {object} SuccessResponse "Email verified successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or verification code"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Failed to verify email"
// @Router  /auth/verify-email [post]
func VerifyEmail(c *gin.Context) {
	var input VerificationRequest

	// Binding JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var user models.User
	// Find user by email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Check if the verification code is correct
	if user.VerificationCode != input.Code {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid verification code"})
		return
	}

	// Update user to set email as verified
	user.EmailVerified = true
	user.VerificationCode = "" // Optionally clear the verification code
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email verified successfully!",
	})
}

// Login handles user authentication
// @Summary User login
// @Description This endpoint allows users to log in by providing email and password. A JWT token will be returned upon successful login.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials  body  LoginCredentials  true  "User credentials (email and password)"
// @Success 200 {object} SuccessResponse "JWT token"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Unauthorized, invalid credentials or email not verified"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Error generating token or database error"
// @Router  /auth/login [post]
func Login(c *gin.Context) {
	var credentials LoginCredentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", credentials.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error"})
		}
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Email not verified. Please verify your email first."})
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid password"})
		return
	}

	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Login successful",
		Data:    gin.H{"token": tokenString},
	})
}
