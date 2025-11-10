package utils

import (
	"go-crud-boilerplate/models"

	"github.com/labstack/echo/v4"
)

func SendResponse(c echo.Context, data any, message string, status int) error {
	response := models.APIResponse{
		Data:    data,
		Message: message,
	}
	return c.JSON(status, response)
}
