package car

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *carService) GetCars(ctx context.Context, filter *models.CarFilters) ([]*models.CarShort, error) {
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
		orderBy := models.CarOrderByCreatedAt
		filter.OrderBy = &orderBy
	}
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		filter.Status = &models.CarStatusAvailable
	}
	return s.repo.GetCarsByFilter(ctx, filter)
}
