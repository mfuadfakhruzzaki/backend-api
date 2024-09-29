// main.go
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/models"
	"github.com/mfuadfakhruzzaki/backend-api/routes"
	"github.com/mfuadfakhruzzaki/backend-api/seeds"
)

func main() {
	// Menghubungkan ke database
	config.ConnectDatabase()

	// Melakukan migrasi model
	err := config.DB.AutoMigrate(&models.User{}, &models.Package{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// Menjalankan seeding data paket
	seeds.SeedPackages()

	// Membuat router baru dengan Gin
	router := gin.Default()

	// Menambahkan middleware untuk logging dan recovery (optional but recommended)
	// Gin's Default() already includes Logger and Recovery middleware
	// If you used gin.New(), you would need to add them manually
	// router := gin.New()
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// Mendaftarkan semua route API terlebih dahulu
	routes.RegisterRoutes(router)

	// Menambahkan log untuk semua route yang terdaftar
	logRoutes(router)

	// Menyajikan file statis
	router.Static("/uploads", "./uploads")

	// Menjalankan server pada port 8080
	fmt.Println("Server berjalan pada port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// logRoutes logs all registered routes.
// Gin does not provide a built-in way to iterate through routes,
// so we need to access the underlying router's tree.
func logRoutes(router *gin.Engine) {
	for _, route := range router.Routes() {
		fmt.Printf("Route registered: %s %s\n", route.Method, route.Path)
	}
}
