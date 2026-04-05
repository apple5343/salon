package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *carService) UpdateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, nil, nil, ErrForbidden
	}
	if err := c.BeforeUpdate(s.clock); err != nil {
		return nil, nil, nil, nil, errorx.Wrap("update car", errorx.BadRequest, err)
	}
	car, _, _, _, err := s.getCarByID(ctx, c.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if car.Status == models.CarStatusSold || car.Status == models.CarStatusBooked {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to change status directly.", errorx.BadRequest)
	}
	if c.Status == models.CarStatusSold || c.Status == models.CarStatusBooked {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to change the sales status directly.", errorx.BadRequest)
	}
	c, err = s.repo.UpdateCar(ctx, c)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, nil, nil, ErrCarNotFound
		} else if errors.Is(err, repo.ErrForeignKey) {
			return nil, nil, nil, nil, ErrForeignKey
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, nil, nil, nil, ErrCarExists
		}
		return nil, nil, nil, nil, errorx.Wrap("update car", errorx.Internal, err)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeCar,
		EntityID:   c.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.CarPayload(c),
		CreatedAt:  s.clock.Now(),
	})
	return s.getCarByID(ctx, c.ID)
}
