package main

import (
	"context"
	"go-crud-boilerplate/internal/config"
	"go-crud-boilerplate/internal/handler"
	"go-crud-boilerplate/internal/repository"
	"go-crud-boilerplate/internal/router"
	"go-crud-boilerplate/internal/service"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.Load()

	pool, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	if err := config.RunMigrations(cfg); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	config.SeedRolesAndAdmin(pool, cfg)

	userRepo := repository.NewUserRepo(pool)
	roleRepo := repository.NewRoleRepo(pool)

	userSvc := service.NewUserService(userRepo, roleRepo)
	roleSvc := service.NewRoleService(roleRepo)

	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	router.Register(e, userHandler, roleHandler)

	sc := echo.StartConfig{Address: ":" + cfg.ServerPort}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
