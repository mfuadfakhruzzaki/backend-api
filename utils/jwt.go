// utils/jwt.go
package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

// Pastikan JWT_SECRET ter-set
func init() {
    if JWT_SECRET == "" {
        fmt.Println("JWT_SECRET is not set in the environment")
    }
}

// GenerateJWT membuat token JWT berdasarkan email pengguna
func GenerateJWT(email string) (string, error) {
    mySigningKey := []byte(JWT_SECRET)

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)

    claims["authorized"] = true
    claims["email"] = email
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        return "", err
    }

    fmt.Println("Token generated:", tokenString) // Debugging output
    return tokenString, nil
}

// ValidateToken memvalidasi token JWT dan mengembalikan email pengguna jika valid
func ValidateToken(tokenString string) (string, error) {
    mySigningKey := []byte(JWT_SECRET)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return mySigningKey, nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        email, ok := claims["email"].(string)
        if !ok {
            return "", errors.New("invalid token claims")
        }
        return email, nil
    }

    return "", errors.New("invalid token")
}
