// controllers/oauthController.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"github.com/mfuadfakhruzzaki/backend-api/oauth"
	"github.com/mfuadfakhruzzaki/backend-api/utils"
	"gorm.io/gorm"
)

// GoogleUser represents the structure of the user data returned by Google
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GithubUser represents the structure of the user data returned by GitHub
type GithubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GoogleLogin initiates the OAuth flow with Google
// @Summary Initiate Google OAuth
// @Description Redirects the user to the Google OAuth login page
// @Tags OAuth
// @Success 302 "Redirects to Google OAuth login"
// @Router /auth/google/login [get]
func GoogleLogin(c *gin.Context) {
	url := oauth.GoogleOauthConfig.AuthCodeURL(oauth.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the callback from Google after user authentication
// @Summary Handle Google OAuth callback
// @Description This endpoint handles the callback from Google after the user has authenticated. It logs in the user or creates a new user account if the user does not already exist.
// @Tags OAuth
// @Param state query string true "OAuth State"
// @Param code query string true "OAuth Code"
// @Success 200 {object} map[string]interface{} "JWT Token"
// @Failure 400 {object} map[string]interface{} "Invalid OAuth state or code"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/google/callback [get]
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauth.OauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	content, err := oauth.GetUserDataFromGoogle(code)
	if err != nil {
		// Redirect to home or an error page
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Parse JSON data from Google
	var googleUser GoogleUser
	err = json.Unmarshal(content, &googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user data"})
		return
	}

	// Check if the user already exists in the database
	var user models.User
	result := config.DB.Where("email = ?", googleUser.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// If user doesn't exist, create a new user
		user = models.User{
			Email:          googleUser.Email,
			Username:       googleUser.Name,
			Password:       "", // Password can be empty or set a default value since OAuth is used
			ProfilePicture: googleUser.Picture,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Return the token to the user
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// GithubLogin initiates the OAuth flow with GitHub
// @Summary Initiate GitHub OAuth
// @Description Redirects the user to the GitHub OAuth login page
// @Tags OAuth
// @Success 302 "Redirects to GitHub OAuth login"
// @Router /auth/github/login [get]
func GithubLogin(c *gin.Context) {
	url := oauth.GithubOauthConfig.AuthCodeURL(oauth.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GithubCallback handles the callback from GitHub after user authentication
// @Summary Handle GitHub OAuth callback
// @Description This endpoint handles the callback from GitHub after the user has authenticated. It logs in the user or creates a new user account if the user does not already exist.
// @Tags OAuth
// @Param state query string true "OAuth State"
// @Param code query string true "OAuth Code"
// @Success 200 {object} map[string]interface{} "JWT Token"
// @Failure 400 {object} map[string]interface{} "Invalid OAuth state or code"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/github/callback [get]
func GithubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauth.OauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	content, err := oauth.GetUserDataFromGithub(code)
	if err != nil {
		// Redirect to home or an error page
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Parse JSON data from GitHub
	var githubUser GithubUser
	err = json.Unmarshal(content, &githubUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user data"})
		return
	}

	// Ensure email is available
	email := githubUser.Email
	if email == "" {
		// Attempt to fetch email via another API call if necessary
		// Alternatively, prompt the user to provide an email
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not available from GitHub"})
		return
	}

	// Check if the user already exists in the database
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// If user doesn't exist, create a new user
		user = models.User{
			Email:          email,
			Username:       githubUser.Login,
			Password:       "", // Password can be empty or set a default value since OAuth is used
			ProfilePicture: githubUser.AvatarURL,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Return the token to the user
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
