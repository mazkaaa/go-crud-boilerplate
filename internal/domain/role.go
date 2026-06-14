package domain

import (
	"context"
	"time"
)

type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Users     []User    `json:"users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleRepository interface {
	FindAll(ctx context.Context) ([]Role, error)
	Create(ctx context.Context, name string) (Role, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Role, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	FindByIDWithUsers(ctx context.Context, id string) (Role, error)
}
