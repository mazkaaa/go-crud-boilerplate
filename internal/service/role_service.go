package service

import (
	"context"
	"errors"
	"go-crud-boilerplate/internal/domain"
)

type RoleService struct {
	repo domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetRoles(ctx context.Context, params domain.PaginationParams) (domain.PaginatedResult, error) {
	params.Sanitize()
	result, err := s.repo.FindAllPaginated(ctx, params)
	if err != nil {
		return domain.PaginatedResult{}, domain.ErrInternal
	}
	return result, nil
}

func (s *RoleService) CreateRole(ctx context.Context, name string) (domain.Role, error) {
	if name == "" {
		return domain.Role{}, domain.ErrInvalidInput
	}

	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return domain.Role{}, domain.ErrInternal
	}
	if exists {
		return domain.Role{}, domain.ErrConflict
	}

	role, err := s.repo.Create(ctx, name)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Role{}, domain.ErrConflict
		}
		return domain.Role{}, domain.ErrInternal
	}
	return role, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return domain.ErrInternal
	}
	if !exists {
		return domain.ErrNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (s *RoleService) GetDetailRole(ctx context.Context, id string) (domain.Role, error) {
	role, err := s.repo.FindByIDWithUsers(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Role{}, domain.ErrNotFound
		}
		return domain.Role{}, domain.ErrInternal
	}
	return role, nil
}
