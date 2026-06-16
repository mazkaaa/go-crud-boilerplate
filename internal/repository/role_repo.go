package repository

import (
	"context"
	"fmt"
	"go-crud-boilerplate/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	pool *pgxpool.Pool
}

func NewRoleRepo(pool *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{pool: pool}
}

func (r *RoleRepo) FindAll(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, created_at, updated_at FROM roles")
	if err != nil {
		return nil, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var ro domain.Role
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.CreatedAt, &ro.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, ro)
	}
	return roles, nil
}

var roleAllowedSorts = map[string]string{
	"id":         "id",
	"name":       "name",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (r *RoleRepo) FindAllPaginated(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult, error) {
	sortBy := sanitizeSortBy(params.SortBy, roleAllowedSorts)
	offset := (params.Page - 1) * params.Limit

	query := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM roles ORDER BY %s %s LIMIT $1 OFFSET $2",
		sortBy, params.SortOrder)

	rows, err := r.pool.Query(ctx, query, params.Limit, offset)
	if err != nil {
		return domain.PaginatedResult{}, fmt.Errorf("query roles paginated: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var ro domain.Role
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.CreatedAt, &ro.UpdatedAt); err != nil {
			return domain.PaginatedResult{}, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, ro)
	}

	var total int
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM roles").Scan(&total)
	if err != nil {
		return domain.PaginatedResult{}, fmt.Errorf("count roles: %w", err)
	}

	return domain.PaginatedResult{
		Items:      roles,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: domain.ComputeTotalPages(total, params.Limit),
	}, nil
}

func (r *RoleRepo) Create(ctx context.Context, name string) (domain.Role, error) {
	var ro domain.Role
	err := r.pool.QueryRow(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at, updated_at", name,
	).Scan(&ro.ID, &ro.Name, &ro.CreatedAt, &ro.UpdatedAt)
	if err != nil {
		return ro, fmt.Errorf("insert role: %w", err)
	}
	return ro, nil
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

func (r *RoleRepo) FindByID(ctx context.Context, id string) (domain.Role, error) {
	var ro domain.Role
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, created_at, updated_at FROM roles WHERE id = $1", id,
	).Scan(&ro.ID, &ro.Name, &ro.CreatedAt, &ro.UpdatedAt)
	if err != nil {
		return ro, fmt.Errorf("find role by id: %w", err)
	}
	return ro, nil
}

func (r *RoleRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check role exists by name: %w", err)
	}
	return exists, nil
}

func (r *RoleRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check role exists by id: %w", err)
	}
	return exists, nil
}

func (r *RoleRepo) FindByIDWithUsers(ctx context.Context, id string) (domain.Role, error) {
	ro, err := r.FindByID(ctx, id)
	if err != nil {
		return ro, err
	}

	rows, err := r.pool.Query(ctx,
		"SELECT id, name, email, role_id, created_at, updated_at FROM users WHERE role_id = $1", id)
	if err != nil {
		return ro, fmt.Errorf("query users for role: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return ro, fmt.Errorf("scan user for role: %w", err)
		}
		ro.Users = append(ro.Users, u)
	}
	return ro, nil
}
