package brand

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *brandService) Delete(ctx context.Context, id string) error {
	role := ctxutil.UserRoleFromContext(ctx)
	if role != string(models.EmployeeRoleAdmin) {
		return ErrForbidden
	}
	userID := ctxutil.UserIDFromContext(ctx)
	brand, err := s.getByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return errorx.NewError("delete brand: "+err.Error(), errorx.Internal)
	}
	s.eventService.AddEvent(ctx, &models.Event{
		Type:       models.EventTypeDeleted,
		EntityType: models.EntityTypeBrand,
		EntityID:   id,
		ActorID:    userID,
		ActorRole:  role,
		Payload:    models.BrandPayload(brand),
	})
	return nil
}
