package model

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *modelService) Update(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, nil, ErrForbidden
	}
	if err := m.BeforeUpdate(s.clock); err != nil {
		return nil, nil, errorx.Wrap("update model", errorx.BadRequest, err)
	}
	m, err := s.repo.Update(ctx, m)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, nil, ErrModelNotFound
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, nil, ErrModelExists
		} else if errors.Is(err, repo.ErrForeignKey) {
			return nil, nil, ErrForeignKey
		}
		return nil, nil, errorx.Wrap("update model", errorx.Internal, err)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeModel,
		EntityID:   m.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.ModelPayload(m),
		CreatedAt:  s.clock.Now(),
	})
	b, err := s.brandService.GetByID(ctx, m.BrandID)
	if err != nil {
		return nil, nil, errorx.Wrap("update model", errorx.Internal, err)
	}
	return m, b, nil
}
