package handler

import (
	"errors"
	"go-crud-boilerplate/internal/domain"
	"go-crud-boilerplate/internal/dto"
	"go-crud-boilerplate/internal/service"
	"go-crud-boilerplate/pkg/response"
	"net/http"

	"github.com/labstack/echo/v5"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) GetRoles(c *echo.Context) error {
	roles, err := h.svc.GetRoles(c.Request().Context())
	if err != nil {
		return response.Send(c, nil, "Failed to fetch roles", http.StatusInternalServerError)
	}
	return response.Send(c, roles, "Success", http.StatusOK)
}

func (h *RoleHandler) CreateRole(c *echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Send(c, nil, "Invalid request body", http.StatusBadRequest)
	}

	role, err := h.svc.CreateRole(c.Request().Context(), req.Name)
	if err != nil {
		return handleRoleError(c, err)
	}
	return response.Send(c, role, "Role created successfully", http.StatusCreated)
}

func (h *RoleHandler) DeleteRole(c *echo.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteRole(c.Request().Context(), id); err != nil {
		return handleRoleError(c, err)
	}
	return response.Send(c, nil, "Role deleted successfully", http.StatusOK)
}

func (h *RoleHandler) GetDetailRole(c *echo.Context) error {
	id := c.Param("id")
	role, err := h.svc.GetDetailRole(c.Request().Context(), id)
	if err != nil {
		return handleRoleError(c, err)
	}
	return response.Send(c, role, "Success", http.StatusOK)
}

func handleRoleError(c *echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return response.Send(c, nil, "Invalid input", http.StatusBadRequest)
	case errors.Is(err, domain.ErrNotFound):
		return response.Send(c, nil, "Resource not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrConflict):
		return response.Send(c, nil, "Resource already exists", http.StatusConflict)
	default:
		return response.Send(c, nil, "Internal server error", http.StatusInternalServerError)
	}
}
