// config/config.go
package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    dsn := "host=localhost user=postgres password=yourpassword dbname=backend_api_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("Gagal terhubung ke database!")
    }

    DB = database
    fmt.Println("Database terhubung!")
}
