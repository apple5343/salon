package employee

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"

	"github.com/apple5343/errorx"
)

func (s *employeeService) getByID(ctx context.Context, id string) (*models.Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, errorx.NewError("get employee: "+err.Error(), errorx.Internal)
	}
	return employee, nil
}

func (s *employeeService) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	//TODO инофрмацию может получать только сам сотрудник или админ
	return s.getByID(ctx, id)
}
