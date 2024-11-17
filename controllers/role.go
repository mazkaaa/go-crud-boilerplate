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
