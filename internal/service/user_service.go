package service

import (
	"context"
	"errors"
	"go-crud-boilerplate/internal/domain"
	"go-crud-boilerplate/pkg/hash"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserService struct {
	repo domain.UserRepository
	roleRepo domain.RoleRepository
}

func NewUserService(repo domain.UserRepository, roleRepo domain.RoleRepository) *UserService {
	return &UserService{repo: repo, roleRepo: roleRepo}
}

func (s *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, domain.ErrInternal
	}
	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password, roleID string) (domain.User, error) {
	if name == "" || email == "" || password == "" || roleID == "" {
		return domain.User{}, domain.ErrInvalidInput
	}

	exists, err := s.roleRepo.ExistsByID(ctx, roleID)
	if err != nil {
		return domain.User{}, domain.ErrInternal
	}
	if !exists {
		return domain.User{}, domain.ErrNotFound
	}

	hashedPassword, err := hash.Password(password)
	if err != nil {
		return domain.User{}, domain.ErrInternal
	}

	user, err := s.repo.Create(ctx, name, email, hashedPassword, roleID)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrConflict
		}
		return domain.User{}, domain.ErrInternal
	}
	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
