package car

import (
	"context"
	"salon/internal/models"

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
	return s.repo.GetCarsByFilter(ctx, filter)
}
