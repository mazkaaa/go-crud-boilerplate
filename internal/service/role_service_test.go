package service

import (
	"context"
	"go-crud-boilerplate/internal/domain"
	"testing"
	"time"
)

func TestGetRoles(t *testing.T) {
	repo := &mockRoleRepo{
		roles: []domain.Role{
			{ID: "1", Name: "admin"},
			{ID: "2", Name: "user"},
		},
	}
	svc := NewRoleService(repo)

	result, err := svc.GetRoles(context.Background(), domain.PaginationParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roles := result.Items.([]domain.Role)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}
}

func TestCreateRoleSuccess(t *testing.T) {
	repo := &mockRoleRepo{
		exists: false,
		role:   domain.Role{ID: "1", Name: "editor"},
	}
	svc := NewRoleService(repo)

	role, err := svc.CreateRole(context.Background(), "editor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Name != "editor" {
		t.Fatalf("expected name editor, got %s", role.Name)
	}
}

func TestCreateRoleDuplicate(t *testing.T) {
	repo := &mockRoleRepo{exists: true}
	svc := NewRoleService(repo)

	_, err := svc.CreateRole(context.Background(), "admin")
	if err != domain.ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestCreateRoleEmptyName(t *testing.T) {
	svc := NewRoleService(&mockRoleRepo{})

	_, err := svc.CreateRole(context.Background(), "")
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestDeleteRoleSuccess(t *testing.T) {
	repo := &mockRoleRepo{exists: true}
	svc := NewRoleService(repo)

	err := svc.DeleteRole(context.Background(), "role-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteRoleNotFound(t *testing.T) {
	repo := &mockRoleRepo{exists: false}
	svc := NewRoleService(repo)

	err := svc.DeleteRole(context.Background(), "bad-id")
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetDetailRole(t *testing.T) {
	repo := &mockRoleRepo{
		role: domain.Role{
			ID:        "1",
			Name:      "admin",
			Users:     []domain.User{{ID: "u1", Name: "Alice", Email: "alice@example.com"}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	svc := NewRoleService(repo)

	role, err := svc.GetDetailRole(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(role.Users) != 1 {
		t.Fatalf("expected 1 user in role, got %d", len(role.Users))
	}
}
