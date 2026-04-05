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
	e.Status = models.EmployeeStatusInactive
	if err := e.BeforeCreate(s.clock); err != nil {
		return nil, errorx.Wrap("register employee", errorx.BadRequest, err)
	}
	e.HireDate = e.CreatedAt
	e, err := s.repo.Create(ctx, e)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("employee alredy exists", errorx.Conflict)
		}
		return nil, errorx.Wrap("register employee", errorx.Internal, err)
	}
	//TODO добавить логи
	return e, nil
}

func (s *employeeService) Hire(ctx context.Context, id string) (*models.Employee, error) {
	userRole := ctxutil.UserRoleFromContext(ctx)
	if userRole == "" {
		return nil, ErrUnauthorized
	}
	if userRole != string(models.EmployeeRoleAdmin) {
		return nil, errorx.NewError("only admin can hire employee", errorx.Forbidden)
	}
	e, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.Status != models.EmployeeStatusInactive {
		return nil, errorx.NewError("employee already hired", errorx.BadRequest)
	}
	e.Status = models.EmployeeStatusActive
	e.HireDate = s.clock.Now()
	e, err = s.repo.Update(ctx, e)
	if err != nil {
		return nil, errorx.Wrap("hire employee", errorx.Internal, err)
	}
	//TODO добавить логи
	return e, nil
}

// from cli only
func (s *employeeService) RegisterAdmin(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	if err := e.BeforeCreate(s.clock); err != nil {
		return nil, errorx.Wrap("register employee", errorx.BadRequest, err)
	}
	e.HireDate = e.CreatedAt
	e.Status = models.EmployeeStatusActive
	e, err := s.repo.Create(ctx, e)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("employee alredy exists", errorx.Conflict)
		}
		return nil, errorx.Wrap("register employee", errorx.Internal, err)
	}
	//TODO добавить логи
	return e, nil
}
