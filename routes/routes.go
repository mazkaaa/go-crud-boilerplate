package routes

import (
	"go-crud-boilerplate/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/users", controllers.GetUsers)
	e.GET("/roles", controllers.GetRoles)
}
