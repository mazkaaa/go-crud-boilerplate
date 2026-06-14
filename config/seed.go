package config

import (
	"context"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func SeedRolesAndAdmin() {
	var count int
	err := Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM roles").Scan(&count)
	if err != nil {
		log.Fatalf("failed to count roles: %v", err)
	}

	if count == 0 {
		var roleID string
		err := Pool.QueryRow(context.Background(),
			"INSERT INTO roles (name) VALUES ($1) RETURNING id", "admin",
		).Scan(&roleID)
		if err != nil {
			log.Fatalf("failed to create default role: %v", err)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SEED_PASSWORD")), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		_, err = Pool.Exec(context.Background(),
			"INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, $4)",
			os.Getenv("SEED_NAME"),
			os.Getenv("SEED_EMAIL"),
			string(passwordHash),
			roleID,
		)
		if err != nil {
			log.Fatalf("failed to create default admin user: %v", err)
		}
		log.Println("default admin role and user created successfully")
	}
}
