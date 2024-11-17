package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),     // e.g., "localhost"
		os.Getenv("DB_USER"),     // e.g., "postgres"
		os.Getenv("DB_PASSWORD"), // e.g., "password"
		os.Getenv("DB_NAME"),     // e.g., "mydb"
		os.Getenv("DB_PORT"),     // e.g., "5432"
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	DB = db
	fmt.Println("database connection establised!")
}
