package controllers

import (
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"go-crud-boilerplate/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRoles(c echo.Context) error {
	var roles []models.Role

	if err := config.DB.Find(&roles).Error; err != nil {
		return utils.SendResponse(c, nil, "Failed to fetch data", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, roles, "Success", http.StatusOK)
}

func GetDetailRole(c echo.Context) error {
	id := c.Param("id")
	var role models.Role

	if err := config.DB.Preload("Users").First(&role, "id = ?", id).Error; err != nil {
		return utils.SendResponse(c, nil, "Role not found", http.StatusNotFound)
	}
	return utils.SendResponse(c, role, "Success", http.StatusOK)
}
