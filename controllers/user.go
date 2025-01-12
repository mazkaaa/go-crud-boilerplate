package controllers

import (
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"go-crud-boilerplate/utils"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		return utils.SendResponse(c, nil, "Failed to fetch data", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, users, "Success", http.StatusOK)
}

func CreateUser(c echo.Context) error {

	findRole := config.DB.Where("id = ?", c.FormValue("role_id")).First(&models.Role{})
	if findRole.Error != nil {
		return utils.SendResponse(c, nil, "Role not found", http.StatusNotFound)
	}

	hashedPassword := utils.HashPassword(c.FormValue("password"))

	user := &models.User{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: hashedPassword,
		RoleID:   c.FormValue("role_id"),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Fatalf("failed to create user: %v", err)
		return utils.SendResponse(c, nil, err.Error(), http.StatusInternalServerError)
	}
	log.Println("User created successfully with ID: ", user.ID)
	return utils.SendResponse(c, user, "User created successfully", http.StatusCreated)
}
