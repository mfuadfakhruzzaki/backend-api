// oauth/oauth.go
package oauth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
    GoogleOauthConfig *oauth2.Config
    GithubOauthConfig *oauth2.Config
    OauthStateString  = "random" // Sebaiknya gunakan string acak yang aman
)

func init() {
    // Memuat file .env
    err := godotenv.Load()
    if err != nil {
        // Jika gagal, gunakan variabel lingkungan
        log.Println("Gagal memuat file .env, menggunakan variabel lingkungan")
    }

    GoogleOauthConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:8080/auth/google/callback",
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
        Endpoint:     google.Endpoint,
    }

    GithubOauthConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:8080/auth/github/callback",
        ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
        ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
        Scopes:       []string{"user:email"},
        Endpoint:     github.Endpoint,
    }
}

func GetUserDataFromGoogle(code string) ([]byte, error) {
    token, err := GoogleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        return nil, err
    }

    response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()
    return ioutil.ReadAll(response.Body)
}

func GetUserDataFromGithub(code string) ([]byte, error) {
    token, err := GithubOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        return nil, err
    }

    // Membuat klien HTTP dengan token
    client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

    // Mendapatkan data pengguna
    userResponse, err := client.Get("https://api.github.com/user")
    if err != nil {
        return nil, err
    }
    defer userResponse.Body.Close()
    userData, err := ioutil.ReadAll(userResponse.Body)
    if err != nil {
        return nil, err
    }

    // Mendapatkan email pengguna
    emailResponse, err := client.Get("https://api.github.com/user/emails")
    if err != nil {
        return nil, err
    }
    defer emailResponse.Body.Close()
    emailData, err := ioutil.ReadAll(emailResponse.Body)
    if err != nil {
        return nil, err
    }

    // Gabungkan data user dan email menjadi satu map
    var userMap map[string]interface{}
    var emailArray []map[string]interface{}

    err = json.Unmarshal(userData, &userMap)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(emailData, &emailArray)
    if err != nil {
        return nil, err
    }

    // Tambahkan email ke userMap
    if len(emailArray) > 0 {
        for _, emailObj := range emailArray {
            if primary, ok := emailObj["primary"].(bool); ok && primary {
                if email, ok := emailObj["email"].(string); ok {
                    userMap["email"] = email
                    break
                }
            }
        }
    }

    // Kembalikan data pengguna sebagai JSON
    return json.Marshal(userMap)
}
