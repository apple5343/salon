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

func (s *supplierService) Delete(ctx context.Context, id string) error {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return ErrForbidden
	}
	supplier, err := s.getByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrForeignKey) {
			return errorx.NewError("other entities depend on this supplier", errorx.BadRequest)
		}
		return errorx.Wrap("delete supplier", errorx.Internal, err)
	}
	if err = s.cache.DeleteByID(ctx, id); err != nil && !errors.Is(err, repo.ErrNotFound) {
		logger.FromContextOrDefault(ctx).Error(ctx, "delete supplier: "+err.Error())
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeDeleted,
		EntityType: models.EntityTypeSupplier,
		EntityID:   supplier.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.SupplierPayload(supplier),
		CreatedAt:  s.clock.Now(),
	})
	return nil
}
