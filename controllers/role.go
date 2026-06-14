package controllers

import (
	"context"
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"go-crud-boilerplate/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetRoles(c *echo.Context) error {
	rows, err := config.Pool.Query(context.Background(),
		"SELECT id, name, created_at, updated_at FROM roles")
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to fetch data", http.StatusInternalServerError)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return utils.SendResponse(c, nil, "Failed to scan role", http.StatusInternalServerError)
		}
		roles = append(roles, r)
	}

	return utils.SendResponse(c, roles, "Success", http.StatusOK)
}

func CreateRole(c *echo.Context) error {
	name := c.FormValue("name")

	var exists bool
	err := config.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to check role", http.StatusInternalServerError)
	}
	if exists {
		return utils.SendResponse(c, nil, "Role already exists", http.StatusConflict)
	}

	var role models.Role
	err = config.Pool.QueryRow(context.Background(),
		"INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at, updated_at", name,
	).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return utils.SendResponse(c, nil, "Failed to create role", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, role, "Role created successfully", http.StatusCreated)
}

func DeleteRole(c *echo.Context) error {
	id := c.Param("id")

	var exists bool
	err := config.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to check role", http.StatusInternalServerError)
	}
	if !exists {
		return utils.SendResponse(c, nil, "Role doesn't exist", http.StatusNotFound)
	}

	_, err = config.Pool.Exec(context.Background(), "DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to delete role", http.StatusInternalServerError)
	}

	return utils.SendResponse(c, nil, "Role deleted successfully", http.StatusOK)
}

func GetDetailRole(c *echo.Context) error {
	id := c.Param("id")

	var role models.Role
	err := config.Pool.QueryRow(context.Background(),
		"SELECT id, name, created_at, updated_at FROM roles WHERE id = $1", id,
	).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return utils.SendResponse(c, nil, "Role not found", http.StatusNotFound)
	}

	rows, err := config.Pool.Query(context.Background(),
		"SELECT id, name, email, role_id, created_at, updated_at FROM users WHERE role_id = $1", id)
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to fetch users", http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return utils.SendResponse(c, nil, "Failed to scan user", http.StatusInternalServerError)
		}
		role.Users = append(role.Users, u)
	}

	return utils.SendResponse(c, role, "Success", http.StatusOK)
}
