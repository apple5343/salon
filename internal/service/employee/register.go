package employee

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *employeeService) Register(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	userRole := ctxutil.UserRoleFromContext(ctx)
	if userRole == "" {
		return nil, ErrUnauthorized
	}
	if userRole != string(models.EmployeeRoleAdmin) {
		return nil, errorx.NewError("only admin can register employee", errorx.Forbidden)
	}
	if e.Role == models.EmployeeRoleAdmin {
		return nil, errorx.NewError("admin cannot be registered", errorx.BadRequest)
	}

	if err := e.BeforeCreate(s.clock); err != nil {
		return nil, errorx.NewError("register employee: "+err.Error(), errorx.BadRequest)
	}
	e.HireDate = e.CreatedAt
	e, err := s.repo.Create(ctx, e)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("employee alredy exists", errorx.Conflict)
		}
		return nil, errorx.NewError("register employee: "+err.Error(), errorx.Internal)
	}
	//TODO добавить логи
	return e, nil
}

// from cli only
func (s *employeeService) RegisterAdmin(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if err := e.BeforeCreate(s.clock); err != nil {
		return nil, errorx.NewError("register employee: "+err.Error(), errorx.BadRequest)
	}
	e.HireDate = e.CreatedAt
	e, err := s.repo.Create(ctx, e)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("employee alredy exists", errorx.Conflict)
		}
		return nil, errorx.NewError("register employee: "+err.Error(), errorx.Internal)
	}
	//TODO добавить логи
	return e, nil
}
