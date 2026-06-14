package controllers

import (
	"context"
	"go-crud-boilerplate/config"
	"go-crud-boilerplate/models"
	"go-crud-boilerplate/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetUsers(c *echo.Context) error {
	rows, err := config.Pool.Query(context.Background(),
		"SELECT id, name, email, role_id, created_at, updated_at FROM users")
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to fetch data", http.StatusInternalServerError)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return utils.SendResponse(c, nil, "Failed to scan user", http.StatusInternalServerError)
		}
		users = append(users, u)
	}

	return utils.SendResponse(c, users, "Success", http.StatusOK)
}

func CreateUser(c *echo.Context) error {
	roleID := c.FormValue("role_id")

	var exists bool
	err := config.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)", roleID).Scan(&exists)
	if err != nil || !exists {
		return utils.SendResponse(c, nil, "Role not found", http.StatusNotFound)
	}

	hashedPassword, err := utils.HashPassword(c.FormValue("password"))
	if err != nil {
		return utils.SendResponse(c, nil, "Failed to hash password", http.StatusInternalServerError)
	}

	var user models.User
	err = config.Pool.QueryRow(context.Background(),
		"INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id, name, email, role_id, created_at, updated_at",
		c.FormValue("name"), c.FormValue("email"), hashedPassword, roleID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return utils.SendResponse(c, nil, err.Error(), http.StatusInternalServerError)
	}

	return utils.SendResponse(c, user, "User created successfully", http.StatusCreated)
}
