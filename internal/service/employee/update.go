package employee

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	"salon/internal/utils/password"

	"github.com/apple5343/errorx"
)

func (s *employeeService) Update(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	existing, err := s.getByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	if existing.Role == models.EmployeeRoleAdmin {
		return nil, errorx.NewError("admin cannot be updated", errorx.BadRequest)
	}
	if existing.Status != e.Status {
		return nil, errorx.NewError("status cannot be changed", errorx.BadRequest)
	}
	if err := e.BeforeUpdate(s.clock); err != nil {
		return nil, errorx.Wrap("update employee", errorx.BadRequest, err)
	}
	if e.PasswordHash != "" {
		hashed, err := password.HashPassword(e.PasswordHash)
		if err != nil {
			return nil, errorx.Wrap("update employee", errorx.BadRequest, err)
		}
		e.PasswordHash = hashed
	} else {
		e.PasswordHash = existing.PasswordHash
	}

	updated, err := s.repo.Update(ctx, e)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("employee alredy exists", errorx.Conflict)
		}
		return nil, errorx.Wrap("update employee", errorx.Internal, err)
	}
	//TODO добавить логи
	return updated, nil
}
