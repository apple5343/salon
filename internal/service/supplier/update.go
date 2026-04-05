package supplier

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	"salon/pkg/logger"

	"github.com/apple5343/errorx"
)

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := supplier.BeforeUpdate(s.clock); err != nil {
		return nil, errorx.Wrap("update supplier", errorx.BadRequest, err)
	}
	supplier, err := s.repo.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrSupplierNotFound
		}
		return nil, errorx.Wrap("update supplier", errorx.Internal, err)
	}
	if err = s.cache.SetByID(ctx, supplier, ttl); err != nil {
		logger.FromContextOrDefault(ctx).Error(ctx, "update supplier: "+err.Error())
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeSupplier,
		EntityID:   supplier.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.SupplierPayload(supplier),
		CreatedAt:  s.clock.Now(),
	})
	return supplier, nil
}
