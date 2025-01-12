package config

import (
	"go-crud-boilerplate/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func SeedRolesAndAdmin() {
	var roleCount int64
	DB.Model(&models.Role{}).Count(&roleCount)

	if roleCount == 0 {
		adminRole := models.Role{Name: "admin"}
		if err := DB.Create(&adminRole).Error; err != nil {
			log.Fatalf("failed to create default role: %v", err)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SEED_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		adminUser := models.User{
			Name:     os.Getenv("SEED_NAME"),
			Email:    os.Getenv("SEED_EMAIL"),
			Password: string(passwordHash),
			Role:     adminRole,
		}

		if err := DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("failed to create default admin user: %v", err)
		}
		log.Println("default admin role and user created successfully")
	}
}
