package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	ServerPort    string
	SeedName      string
	SeedEmail     string
	SeedPassword  string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	return Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBPort:       os.Getenv("DB_PORT"),
		ServerPort:   getEnvOrDefault("SERVER_PORT", "1323"),
		SeedName:     os.Getenv("SEED_NAME"),
		SeedEmail:    os.Getenv("SEED_EMAIL"),
		SeedPassword: os.Getenv("SEED_PASSWORD"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func ConnectDB(cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("database connection established!")
	return pool, nil
}

func RunMigrations(cfg Config) error {
	m, err := migrate.New("file://migrations", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("migrations completed")
	return nil
}

func SeedRolesAndAdmin(pool *pgxpool.Pool, cfg Config) {
	ctx := context.Background()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM roles").Scan(&count)
	if err != nil {
		log.Fatalf("failed to count roles: %v", err)
	}
	if count > 0 {
		return
	}

	var roleID string
	err = pool.QueryRow(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id", "admin",
	).Scan(&roleID)
	if err != nil {
		log.Fatalf("failed to create default role: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	_, err = pool.Exec(ctx,
		"INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, $4)",
		cfg.SeedName, cfg.SeedEmail, string(passwordHash), roleID,
	)
	if err != nil {
		log.Fatalf("failed to create default admin user: %v", err)
	}
	log.Println("default admin role and user created successfully")
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
