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

func (s *carService) CreateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := b.BeforeCreate(); err != nil {
		return nil, errorx.NewError("create brand: "+err.Error(), errorx.BadRequest)
	}
	b, err := s.repo.CreateBrand(ctx, b)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, ErrBrandExists
		}
		return nil, errorx.NewError("create brand: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeCreated,
		EntityType: models.EntityTypeBrand,
		EntityID:   b.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.BrandPayload(b),
		CreatedAt:  time.Now(),
	})
	return b, nil
}

func (s *carService) CreateModel(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, ErrForbidden
	}
	if err := m.BeforeCreate(); err != nil {
		return nil, nil, errorx.NewError("create model: "+err.Error(), errorx.BadRequest)
	}
	m, err := s.repo.CreateModel(ctx, m)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, nil, ErrModelExists
		} else if errors.Is(err, repo.ErrForeignKey) {
			return nil, nil, ErrForeignKey
		}
		return nil, nil, errorx.NewError("create model: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeCreated,
		EntityType: models.EntityTypeModel,
		EntityID:   m.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.ModelPayload(m),
		CreatedAt:  time.Now(),
	})
	b, err := s.getBrandByID(ctx, m.BrandID)
	if err != nil {
		return nil, nil, errorx.NewError("create model: "+err.Error(), errorx.Internal)
	}
	return m, b, nil
}

func (s *carService) CreateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, nil, nil, ErrForbidden
	}
	if err := c.BeforeCreate(); err != nil {
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
		CreatedAt:  time.Now(),
	})
	return s.getCarByID(ctx, c.ID)
}
