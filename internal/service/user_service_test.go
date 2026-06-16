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
		ID:     "generated-id",
		Name:   name,
		Email:  email,
		RoleID: &roleID,
	}
	return u, nil
}

func (m *mockUserRepo) FindByRoleID(ctx context.Context, roleID string) ([]domain.User, error) {
	return m.users, m.err
}

func (m *mockUserRepo) FindAllPaginated(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult, error) {
	if m.err != nil {
		return domain.PaginatedResult{}, m.err
	}
	total := len(m.users)
	return domain.PaginatedResult{
		Items:      m.users,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: domain.ComputeTotalPages(total, params.Limit),
	}, nil
}

type mockRoleRepo struct {
	roles  []domain.Role
	role   domain.Role
	exists bool
	err    error
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

func (m *mockRoleRepo) FindAllPaginated(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult, error) {
	if m.err != nil {
		return domain.PaginatedResult{}, m.err
	}
	total := len(m.roles)
	return domain.PaginatedResult{
		Items:      m.roles,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: domain.ComputeTotalPages(total, params.Limit),
	}, nil
}

func TestGetUsers(t *testing.T) {
	repo := &mockUserRepo{
		users: []domain.User{
			{ID: "1", Name: "Alice", Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
	svc := NewUserService(repo, &mockRoleRepo{})

	result, err := svc.GetUsers(context.Background(), domain.PaginationParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	users := result.Items.([]domain.User)
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}
	if result.Page != 1 {
		t.Fatalf("expected page 1, got %d", result.Page)
	}
	if result.TotalPages != 1 {
		t.Fatalf("expected total_pages 1, got %d", result.TotalPages)
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
