package model

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *modelService) Create(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, ErrForbidden
	}
	if err := m.BeforeCreate(s.clock); err != nil {
		return nil, nil, errorx.NewError("create model: "+err.Error(), errorx.BadRequest)
	}
	m, err := s.repo.Create(ctx, m)
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
		CreatedAt:  s.clock.Now(),
	})
	b, err := s.brandService.GetByID(ctx, m.BrandID)
	if err != nil {
		return nil, nil, errorx.NewError("create model: "+err.Error(), errorx.Internal)
	}
	return m, b, nil
}
