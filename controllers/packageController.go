// controllers/packageController.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/middleware"
	"github.com/mfuadfakhruzzaki/backend-api/models"
)

// GetPackages retrieves all available packages
// @Summary Get all packages
// @Description Retrieve a list of all available packages
// @Tags Packages
// @Produce json
// @Success 200 {array} models.Package "List of available packages"
// @Failure 500 {object} map[string]string "Error fetching packages"
// @Router /packages [get]
func GetPackages(c *gin.Context) {
	var packages []models.Package
	if err := config.DB.Find(&packages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching packages"})
		return
	}

	c.JSON(http.StatusOK, packages)
}

// SelectPackage allows a user to select a package by its ID
// @Summary Select a package
// @Description Allows a user to select a package by its ID, updates the user's selected package
// @Tags Packages
// @Param id path int true "Package ID"
// @Produce json
// @Success 200 {object} map[string]interface{} "Package selected successfully, includes user and package information"
// @Failure 400 {object} map[string]string "Invalid package ID"
// @Failure 401 {object} map[string]string "Unauthorized, user not found in context"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Database error or error updating user package"
// @Router /packages/{id} [post]
func SelectPackage(c *gin.Context) {
	// Retrieve the 'id' parameter from the URL
	packageIDStr := c.Param("id")
	packageID, err := strconv.Atoi(packageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
		return
	}

	// Retrieve the user's email from the context (set by JWT middleware)
	email, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user email in context"})
		return
	}

	// Find the user by email
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

	// Update the user's PackageID
	pkgID := uint(packageID)
	user.PackageID = &pkgID

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user package"})
		return
	}

	// Optionally, you can fetch the updated user or include additional information
	c.JSON(http.StatusOK, gin.H{
		"message":      "Package selected successfully",
		"user":         user,
		"selectedPack": packageID,
	})
}
