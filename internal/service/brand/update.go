package brand

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	"salon/pkg/logger"

	"github.com/apple5343/errorx"
)

func (s *brandService) Update(ctx context.Context, b *models.Brand) (*models.Brand, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := b.BeforeUpdate(s.clock); err != nil {
		return nil, errorx.Wrap("update brand", errorx.BadRequest, err)
	}
	b, err := s.repo.Update(ctx, b)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrBrandNotFound
		} else if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, ErrBrandExists
		}
		return nil, errorx.Wrap("update brand", errorx.Internal, err)
	}
	if err = s.cache.SetByID(ctx, b, ttl); err != nil {
		logger.FromContextOrDefault(ctx).Error(ctx, err.Error())
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeBrand,
		EntityID:   b.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.BrandPayload(b),
		CreatedAt:  s.clock.Now(),
	})
	return b, nil
}
