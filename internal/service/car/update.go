package car

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *carService) UpdateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := b.BeforeUpdate(); err != nil {
		return nil, errorx.NewError("update brand: "+err.Error(), errorx.BadRequest)
	}
	b, err := s.repo.UpdateBrand(ctx, b)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrBrandNotFound
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, ErrBrandExists
		}
		return nil, errorx.NewError("update brand: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeBrand,
		EntityID:   b.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.BrandPayload(b),
		CreatedAt:  time.Now(),
	})
	return b, nil
}

func (s *carService) UpdateModel(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, ErrForbidden
	}
	if err := m.BeforeUpdate(); err != nil {
		return nil, nil, errorx.NewError("update model: "+err.Error(), errorx.BadRequest)
	}
	m, err := s.repo.UpdateModel(ctx, m)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, ErrModelNotFound
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, nil, ErrModelExists
		} else if errors.Is(err, repo.ErrForeignKey) {
			return nil, nil, ErrForeignKey
		}
		return nil, nil, errorx.NewError("update model: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeModel,
		EntityID:   m.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.ModelPayload(m),
		CreatedAt:  time.Now(),
	})
	b, err := s.getBrandByID(ctx, m.BrandID)
	if err != nil {
		return nil, nil, errorx.NewError("update model: "+err.Error(), errorx.Internal)
	}
	return m, b, nil
}

func (s *carService) UpdateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, nil, nil, ErrForbidden
	}
	if err := c.BeforeUpdate(); err != nil {
		return nil, nil, nil, nil, errorx.NewError("update car: "+err.Error(), errorx.BadRequest)
	}
	car, _, _, _, err := s.getCarByID(ctx, c.ID)
	if err != nil {
		return nil, nil, nil, nil, errorx.NewError("update car: "+err.Error(), errorx.Internal)
	}
	if car.Status == models.CarStatusSold {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to set the sales status directly", errorx.BadRequest)
	}
	if c.Status == models.CarStatusSold {
		return nil, nil, nil, nil, errorx.NewError("it is not possible to change the status of a sold vehicle.", errorx.BadRequest)
	}
	c, err = s.repo.UpdateCar(ctx, c)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, nil, nil, ErrCarNotFound
		}
		return nil, nil, nil, nil, errorx.NewError("update car: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeCar,
		EntityID:   c.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.CarPayload(c),
		CreatedAt:  time.Now(),
	})
	return s.getCarByID(ctx, c.ID)
}
