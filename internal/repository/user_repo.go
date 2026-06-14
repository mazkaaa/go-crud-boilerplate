package repository

import (
	"context"
	"fmt"
	"go-crud-boilerplate/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) FindAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, email, role_id, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, name, email, password, roleID string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id, name, email, role_id, created_at, updated_at",
		name, email, password, roleID,
	).Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) FindByRoleID(ctx context.Context, roleID string) ([]domain.User, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, email, role_id, created_at, updated_at FROM users WHERE role_id = $1", roleID)
	if err != nil {
		return nil, fmt.Errorf("query users by role: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}
