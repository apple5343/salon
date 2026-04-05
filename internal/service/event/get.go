package event

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *eventService) getEventByID(ctx context.Context, id string) (*models.Event, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, errorx.Wrap("get event", errorx.Internal, err)
	}
	return e, nil
}

func (s *eventService) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) && role != string(models.EmployeeRoleManager) {
		return nil, ErrForbidden
	}
	return s.getEventByID(ctx, id)
}
