package employee

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *employeeService) GetEmployees(ctx context.Context, filter *models.EmployeeFilters) ([]*models.Employee, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := filter.Validate(); err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}
	if filter.Limit == nil {
		limit := 10
		filter.Limit = &limit
	}
	if filter.Offset == nil {
		offset := 0
		filter.Offset = &offset
	}
	if filter.OrderBy == nil {
		orderBy := models.EmployeeOrderByCreatedAt
		filter.OrderBy = &orderBy
	}
	return s.repo.GetEmployeesByFilter(ctx, filter)
}
