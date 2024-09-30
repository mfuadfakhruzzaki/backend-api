package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/backend-api/config"
	_ "github.com/mfuadfakhruzzaki/backend-api/docs"
	"github.com/mfuadfakhruzzaki/backend-api/routes"
	"github.com/mfuadfakhruzzaki/backend-api/seeds"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Menambahkan rute untuk Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
