package handler

import (
	"errors"
	"go-crud-boilerplate/internal/domain"
	"go-crud-boilerplate/internal/dto"
	"go-crud-boilerplate/internal/service"
	"go-crud-boilerplate/pkg/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUsers(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	params := domain.PaginationParams{
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	result, err := h.svc.GetUsers(c.Request().Context(), params)
	if err != nil {
		return response.Send(c, nil, "Failed to fetch users", http.StatusInternalServerError)
	}

	meta := response.PaginationMeta{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	}
	return response.SendPaginated(c, result.Items, meta, "Success", http.StatusOK)
}

func (h *UserHandler) CreateUser(c *echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Send(c, nil, "Invalid request body", http.StatusBadRequest)
	}

	user, err := h.svc.CreateUser(c.Request().Context(), req.Name, req.Email, req.Password, req.RoleID)
	if err != nil {
		return handleUserError(c, err)
	}
	return response.Send(c, user, "User created successfully", http.StatusCreated)
}

func handleUserError(c *echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return response.Send(c, nil, "Invalid input", http.StatusBadRequest)
	case errors.Is(err, domain.ErrNotFound):
		return response.Send(c, nil, "Role not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrConflict):
		return response.Send(c, nil, "User already exists", http.StatusConflict)
	default:
		return response.Send(c, nil, "Internal server error", http.StatusInternalServerError)
	}
}
