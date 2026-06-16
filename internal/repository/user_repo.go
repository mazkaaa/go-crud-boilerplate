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

var userAllowedSorts = map[string]string{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func sanitizeSortBy(input string, allowed map[string]string) string {
	if col, ok := allowed[input]; ok {
		return col
	}
	return "created_at"
}

func (r *UserRepo) FindAllPaginated(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult, error) {
	sortBy := sanitizeSortBy(params.SortBy, userAllowedSorts)
	offset := (params.Page - 1) * params.Limit

	query := fmt.Sprintf("SELECT id, name, email, role_id, created_at, updated_at FROM users ORDER BY %s %s LIMIT $1 OFFSET $2",
		sortBy, params.SortOrder)

	rows, err := r.pool.Query(ctx, query, params.Limit, offset)
	if err != nil {
		return domain.PaginatedResult{}, fmt.Errorf("query users paginated: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return domain.PaginatedResult{}, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}

	var total int
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return domain.PaginatedResult{}, fmt.Errorf("count users: %w", err)
	}

	return domain.PaginatedResult{
		Items:      users,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: domain.ComputeTotalPages(total, params.Limit),
	}, nil
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
