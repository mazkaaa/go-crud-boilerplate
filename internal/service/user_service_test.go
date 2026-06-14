package service

import (
	"context"
	"go-crud-boilerplate/internal/domain"
	"testing"
)

type mockUserRepo struct {
	users []domain.User
	err   error
}

func (m *mockUserRepo) FindAll(ctx context.Context) ([]domain.User, error) {
	return m.users, m.err
}

func (m *mockUserRepo) Create(ctx context.Context, name, email, password, roleID string) (domain.User, error) {
	if m.err != nil {
		return domain.User{}, m.err
	}
	u := domain.User{
		ID:    "generated-id",
		Name:  name,
		Email: email,
		RoleID: &roleID,
	}
	return u, nil
}

func (m *mockUserRepo) FindByRoleID(ctx context.Context, roleID string) ([]domain.User, error) {
	return m.users, m.err
}

type mockRoleRepo struct {
	roles   []domain.Role
	role    domain.Role
	exists  bool
	err     error
}

func (m *mockRoleRepo) FindAll(ctx context.Context) ([]domain.Role, error) {
	return m.roles, m.err
}

func (m *mockRoleRepo) Create(ctx context.Context, name string) (domain.Role, error) {
	return m.role, m.err
}

func (m *mockRoleRepo) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *mockRoleRepo) FindByID(ctx context.Context, id string) (domain.Role, error) {
	return m.role, m.err
}

func (m *mockRoleRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return m.exists, m.err
}

func (m *mockRoleRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	return m.exists, m.err
}

func (m *mockRoleRepo) FindByIDWithUsers(ctx context.Context, id string) (domain.Role, error) {
	return m.role, m.err
}

func TestGetUsers(t *testing.T) {
	repo := &mockUserRepo{
		users: []domain.User{
			{ID: "1", Name: "Alice", Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
	svc := NewUserService(repo, &mockRoleRepo{})

	users, err := svc.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestCreateUserSuccess(t *testing.T) {
	roleRepo := &mockRoleRepo{exists: true}
	userRepo := &mockUserRepo{}
	svc := NewUserService(userRepo, roleRepo)

	user, err := svc.CreateUser(context.Background(), "Alice", "alice@example.com", "secret", "role-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name Alice, got %s", user.Name)
	}
}

func TestCreateUserRoleNotFound(t *testing.T) {
	roleRepo := &mockRoleRepo{exists: false}
	userRepo := &mockUserRepo{}
	svc := NewUserService(userRepo, roleRepo)

	_, err := svc.CreateUser(context.Background(), "Alice", "alice@example.com", "secret", "bad-role")
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateUserEmptyFields(t *testing.T) {
	roleRepo := &mockRoleRepo{exists: true}
	userRepo := &mockUserRepo{}
	svc := NewUserService(userRepo, roleRepo)

	_, err := svc.CreateUser(context.Background(), "", "email@example.com", "secret", "role-1")
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
