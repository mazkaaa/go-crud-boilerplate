package controllers

import (
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"go-crud-boilerplate/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	hashedPassword := utils.HashPassword(c.FormValue("password"))

	user := &models.User{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: hashedPassword,
		RoleID:   c.FormValue("role_id"),
	}

	emailFind := config.DB.Where("email = ?", user.Email).First(&models.User{})
	roleFind := config.DB.Where("id = ?", user.RoleID).First(&models.Role{})

	if emailFind.Error == nil {
		return utils.SendResponse(c, nil, "Email already exists!", http.StatusBadRequest)
	}
	if roleFind.Error != nil {
		return utils.SendResponse(c, nil, "Role not found!", http.StatusBadRequest)
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return utils.SendResponse(c, nil, err.Error(), http.StatusInternalServerError)
	}
	return utils.SendResponse(c, nil, "Successfully registed", http.StatusCreated)

}
