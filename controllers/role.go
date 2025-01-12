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

func CreateRole(c echo.Context) error {
	role := &models.Role{
		Name: c.FormValue("name"),
	}

	// Check if role already exists
	if err := config.DB.Where("name = ?", role.Name).First(&models.Role{}).Error; err == nil {
		return utils.SendResponse(c, nil, "Role already exists", http.StatusConflict)
	}

	if err := config.DB.Create(&role).Error; err != nil {
		return utils.SendResponse(c, nil, "Failed to create role", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, role, "Role created successfully", http.StatusCreated)
}

func DeleteRole(c echo.Context) error {
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).First(&models.Role{}).Error; err != nil {
		return utils.SendResponse(c, nil, "Role doesn't exists", http.StatusNotFound)
	}

	if err := config.DB.Delete(&models.Role{}, id).Error; err != nil {
		return utils.SendResponse(c, nil, "Failed to delete role", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, nil, "Role deleted successfully", http.StatusOK)
}

func GetDetailRole(c echo.Context) error {
	id := c.Param("id")
	var role models.Role

	if err := config.DB.Preload("Users").First(&role, "id = ?", id).Error; err != nil {
		return utils.SendResponse(c, nil, "Role not found", http.StatusNotFound)
	}
	return utils.SendResponse(c, role, "Success", http.StatusOK)
}
