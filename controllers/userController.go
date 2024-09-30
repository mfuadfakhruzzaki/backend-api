package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/middleware"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"gorm.io/gorm"
)

// UploadProfilePicture handles the upload of a user's profile picture
// @Summary Upload profile picture
// @Description Upload a profile picture for the currently logged-in user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param profile_picture formData file true "Profile picture file (jpg, jpeg, png)"
// @Success 200 {object} map[string]interface{} "Profile picture uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or file type"
// @Failure 401 {object} map[string]interface{} "Unauthorized or email not found"
// @Failure 403 {object} map[string]interface{} "Email not verified"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/profile/picture [post]
func UploadProfilePicture(c *gin.Context) {
	// Log the request for debugging purposes
	fmt.Println("Received request to upload profile picture")

	// Retrieve the user's email from the context (set by JWT middleware)
	email, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Email not found in context"})
		fmt.Println("Unauthorized: Missing or invalid token")
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid email in context"})
		fmt.Println("Unauthorized: Invalid email in context")
		return
	}

	// Find the user in the database based on email
	var user models.User
	result := config.DB.Where("email = ?", emailStr).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			fmt.Printf("User not found: %s\n", emailStr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			fmt.Printf("Database error: %v\n", result.Error)
		}
		return
	}

	// Check if the email is verified before allowing profile picture upload
	if !user.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified. Please verify your email to upload a profile picture."})
		fmt.Printf("User email not verified: %s\n", emailStr)
		return
	}

	// Parse the multipart form with a maximum memory of 10MB
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form data"})
		fmt.Printf("Error parsing form data: %v\n", err)
		return
	}

	// Retrieve the file from the form input named "profile_picture"
	file, handler, err := c.Request.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		fmt.Printf("Error retrieving file: %v\n", err)
		return
	}
	defer file.Close()

	// Validate the file extension
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	fileExt := filepath.Ext(handler.Filename)
	if !allowedExtensions[fileExt] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, JPEG, and PNG are allowed."})
		fmt.Printf("Invalid file type: %s\n", fileExt)
		return
	}

	// Validate the file size (optional, already limited by ParseMultipartForm)
	maxFileSize := int64(10 << 20) // 10MB
	if handler.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB"})
		fmt.Printf("File size exceeds limit: %d bytes\n", handler.Size)
		return
	}

	// Define the directory to save the uploaded files
	uploadDir := "./uploads/profile_pictures/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating upload directory"})
		fmt.Printf("Error creating upload directory: %v\n", err)
		return
	}

	// Create a unique filename using the user's ID and the original file extension
	filename := fmt.Sprintf("user_%d%s", user.ID, fileExt)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file to the specified directory
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		fmt.Printf("Error saving file: %v\n", err)
		return
	}

	// Update the user's ProfilePicture field with the new file path
	user.ProfilePicture = fmt.Sprintf("/uploads/profile_pictures/%s", filename)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user profile"})
		fmt.Printf("Error updating user profile: %v\n", err)
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded successfully", "profile_picture": user.ProfilePicture})
}

// GetProfile returns the profile data of the currently logged-in user
// @Summary Get user profile
// @Description Retrieve the profile of the currently logged-in user
// @Tags User
// @Produce json
// @Success 200 {object} models.User "User profile data"
// @Failure 401 {object} map[string]interface{} "Unauthorized, email not found"
// @Failure 403 {object} map[string]interface{} "Email not verified"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Database error"
// @Router /users/profile [get]
func GetProfile(c *gin.Context) {
	// Retrieve the user's email from the context (set by JWT middleware)
	email, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Email not found in context"})
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid email in context"})
		return
	}

	// Find the user in the database based on email
	var user models.User
	result := config.DB.Where("email = ?", emailStr).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if the email is verified
	if !user.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified. Please verify your email to view profile."})
		return
	}

	// Remove the password field before sending the response for security
	user.Password = ""

	// Return the user's profile data as JSON
	c.JSON(http.StatusOK, user)
}
