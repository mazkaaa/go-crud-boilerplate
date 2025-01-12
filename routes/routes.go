package routes

import (
	"go-crud-boilerplate/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/users", controllers.GetUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/roles", controllers.GetRoles)
	e.POST("/roles", controllers.CreateRole)
	e.GET("/roles/:id", controllers.GetDetailRole)
}
