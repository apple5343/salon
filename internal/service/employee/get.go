package employee

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
	"github.com/docker/distribution/uuid"
)

func (s *employeeService) getByID(ctx context.Context, id string) (*models.Employee, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, errorx.Wrap("get employee", errorx.Internal, err)
	}
	return employee, nil
}

func (s *employeeService) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	return s.getByID(ctx, id)
}

func (s *employeeService) Profile(ctx context.Context) (*models.Employee, error) {
	id := ctxutil.UserIDFromContext(ctx)
	if id == "" {
		return nil, ErrInvalidToken
	}
	return s.getByID(ctx, id)
}
