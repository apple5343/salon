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
	return s.repo.GetCarsByFilter(ctx, filter)
}

func (s *carService) GetModels(ctx context.Context, filter *models.ModelFilters) ([]*models.ModelShort, error) {
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
	return s.repo.GetModelsByFilter(ctx, filter)
}
