package router

import (
	"go-crud-boilerplate/internal/handler"

	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo, uh *handler.UserHandler, rh *handler.RoleHandler) {
	e.GET("/users", uh.GetUsers)
	e.POST("/users", uh.CreateUser)
	e.GET("/roles", rh.GetRoles)
	e.POST("/roles", rh.CreateRole)
	e.GET("/roles/:id", rh.GetDetailRole)
	e.DELETE("/roles/:id", rh.DeleteRole)
}
