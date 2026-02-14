package supplier

import (
	"context"
	"salon/internal/models"

	"github.com/apple5343/errorx"
)

func (s *supplierService) GetSuppliers(ctx context.Context, filters *models.SupplierFilters) ([]*models.Supplier, error) {
	if err := filters.Validate(); err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}
	if filters.Limit == nil {
		limit := 10
		filters.Limit = &limit
	}
	if filters.Offset == nil {
		offset := 0
		filters.Offset = &offset
	}
	if filters.OrderBy == nil {
		order := models.SupplierOrderByCreatedAt
		filters.OrderBy = &order
	}
	return s.repo.GetSuppliersByFilter(ctx, filters)
}
