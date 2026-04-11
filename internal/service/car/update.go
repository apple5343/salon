package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *carService) UpdateCar(ctx context.Context, car *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, nil, nil, ErrForbidden
	}
	if err := car.BeforeUpdate(s.clock); err != nil {
		return nil, nil, nil, nil, errorx.Wrap("update car", errorx.BadRequest, err)
	}
	existingCar, _, _, _, err := s.getCarByID(ctx, car.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if existingCar.Status == models.CarStatusSold || existingCar.Status == models.CarStatusBooked {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to update sold or booked cars.", errorx.BadRequest)
	}
	if car.Status == models.CarStatusSold || car.Status == models.CarStatusBooked {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to change the sales status directly.", errorx.BadRequest)
	}
	car, err = s.repo.UpdateCar(ctx, car)
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
		EntityID:   car.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.CarPayload(car),
		CreatedAt:  s.clock.Now(),
	})
	return s.getCarByID(ctx, car.ID)
}
