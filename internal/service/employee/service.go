package employee

import (
	"salon/internal/config"
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrUnauthorized     = errorx.NewError("unauthorized", errorx.Unauthorized)
	ErrInvalidCreds     = errorx.NewError("invalid credentials", errorx.Unauthorized)
	ErrInvalidToken     = errorx.NewError("invalid token", errorx.Unauthorized)
	ErrEmployeeNotFound = errorx.NewError("employee not found", errorx.BadRequest)
	ErrForbidden        = errorx.NewError("forbidden", errorx.Forbidden)
	ErrInvalidID        = errorx.NewError("invalid id", errorx.BadRequest)
)

type employeeService struct {
	repo      repository.EmployeeRepository
	jwtConfig *config.JWT
}

func NewService(repo repository.EmployeeRepository, tokenConfig *config.JWT) service.EmployeeService {
	return &employeeService{
		repo:      repo,
		jwtConfig: tokenConfig,
	}
}
