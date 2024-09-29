// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    // Memuat file .env
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Gagal memuat file .env, menggunakan variabel lingkungan")
    }

    // Mengambil variabel lingkungan
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
        dbHost, dbUser, dbPassword, dbName, dbPort)
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("Gagal terhubung ke database!")
    }

    DB = database
    fmt.Println("Database terhubung!")
}
