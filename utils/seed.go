package utils

import (
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func SeedRolesAndAdmin() {
	var roleCount int64
	config.DB.Model(&models.Role{}).Count(&roleCount)

	if roleCount == 0 {
		log.Println("No roles found, seeding default roles...")
		SeedDefaultRoles()

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SEED_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		findAdminRole := config.DB.Where("name = ?", "admin").First(&models.Role{})

		adminUser := &models.User{
			Name:     os.Getenv("SEED_NAME"),
			Email:    os.Getenv("SEED_EMAIL"),
			Password: string(passwordHash),
			RoleID:   findAdminRole.Statement.Model.(*models.Role).ID,
		}
		if err := config.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("failed to create admin user: %v", err)
		}

		log.Println("default admin role and user created successfully")
	}
}

func SeedDefaultRoles() {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	for _, role := range roles {
		if err := config.DB.Create(&role).Error; err != nil {
			log.Printf("failed to create role %s: %v", role.Name, err)
		} else {
			log.Printf("role %s created successfully", role.Name)
		}
	}
}
