package supplier

import (
	"context"
	"errors"
	"salon/internal/models"
	repo "salon/internal/repository/errors"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *supplierService) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if err := supplier.BeforeCreate(s.clock); err != nil {
		return nil, errorx.NewError("create supplier: "+err.Error(), errorx.BadRequest)
	}
	supplier, err := s.repo.Create(ctx, supplier)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, errorx.NewError("supplier alredy exists", errorx.Conflict)
		}
		return nil, errorx.NewError("create supplier: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeCreated,
		EntityType: models.EntityTypeSupplier,
		EntityID:   supplier.ID,
		ActorID:    ctxutil.UserIDFromContext(ctx),
		ActorRole:  role,
		Payload:    models.SupplierPayload(supplier),
		CreatedAt:  s.clock.Now(),
	})
	return supplier, nil
}
