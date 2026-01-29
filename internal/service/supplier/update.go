package supplier

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := supplier.BeforeUpdate(); err != nil {
		return nil, errorx.NewError("update supplier: "+err.Error(), errorx.BadRequest)
	}
	supplier, err := s.repo.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrSupplierNotFound
		}
		return nil, errorx.NewError("update supplier: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeUpdated,
		EntityType: models.EntityTypeSupplier,
		EntityID:   supplier.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.SupplierPayload(supplier),
		CreatedAt:  time.Now(),
	})
	return supplier, nil
}
