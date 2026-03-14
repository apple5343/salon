package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *carService) CreateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, nil, nil, ErrForbidden
	}
	if err := c.BeforeCreate(s.clock); err != nil {
		return nil, nil, nil, nil, errorx.NewError("create car: "+err.Error(), errorx.BadRequest)
	}
	c, err := s.repo.CreateCar(ctx, c)
	if err != nil {
		if errors.Is(err, repo.ErrForeignKey) {
			return nil, nil, nil, nil, ErrForeignKey
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, nil, nil, nil, ErrCarExists
		}
		return nil, nil, nil, nil, errorx.NewError("create car: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeCreated,
		EntityType: models.EntityTypeCar,
		EntityID:   c.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.CarPayload(c),
		CreatedAt:  s.clock.Now(),
	})
	return s.getCarByID(ctx, c.ID)
}
