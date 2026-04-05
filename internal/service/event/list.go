package event

import (
	"context"
	"salon/internal/models"

	"github.com/apple5343/errorx"
)

func (s *eventService) GetEvents(ctx context.Context, filters *models.EventFilters) ([]*models.Event, error) {
	if err := filters.Validate(); err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}
	if filters.Limit == nil {
		limit := 10
		filters.Limit = &limit
	}
	if filters.Offset == nil {
		offset := 0
		filters.Offset = &offset
	}

	res, err := s.repo.GetEventsByFilter(ctx, filters)
	if err != nil {
		return nil, errorx.Wrap("get events", errorx.Internal, err)
	}
	return res, nil
}
