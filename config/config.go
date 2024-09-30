package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase menghubungkan ke database dan memuat environment variables
func ConnectDatabase() {
	// Memuat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan variabel lingkungan")
	}

	// Mengambil variabel lingkungan untuk koneksi database
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Data Source Name (DSN) untuk PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Jika gagal terhubung ke database, panic
	if err != nil {
		panic("Gagal terhubung ke database!")
	}

	// Simpan koneksi database ke variabel global
	DB = database
	fmt.Println("Database berhasil terhubung!")

	// Jalankan migrasi untuk menyesuaikan model ke database
	Migrate()
}

// Migrate menjalankan migrasi skema database berdasarkan model yang ada
func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	fmt.Println("Migrasi database berhasil!")
}
