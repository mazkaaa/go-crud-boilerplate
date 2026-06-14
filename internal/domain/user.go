package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	RoleID    *string   `json:"role_id"`
	Role      *Role     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, name, email, password, roleID string) (User, error)
	FindByRoleID(ctx context.Context, roleID string) ([]User, error)
}
