package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/config"
	"github.com/mfuadfakhruzzaki/backend-api/routes"
	"github.com/mfuadfakhruzzaki/backend-api/seeds"
)

func main() {
	// Menghubungkan ke database dan menjalankan migrasi di config.ConnectDatabase()
	config.ConnectDatabase()

	// Menjalankan seeding data paket
	seeds.SeedPackages()

	// Membuat router baru dengan Gin
	router := gin.Default()

	// Mendaftarkan semua route API
	routes.RegisterRoutes(router)

	// Menambahkan log untuk semua route yang terdaftar
	logRoutes(router)

	// Menyajikan file statis dari folder 'uploads'
	router.Static("/uploads", "./uploads")

	// Menjalankan server pada port 8080
	fmt.Println("Server berjalan pada port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// logRoutes logs all registered routes.
func logRoutes(router *gin.Engine) {
	for _, route := range router.Routes() {
		fmt.Printf("Route registered: %s %s\n", route.Method, route.Path)
	}
}
