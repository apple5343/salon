package sale

import (
	"context"
	"salon/internal/models"

	"github.com/apple5343/errorx"
)

func (s *saleService) GetSales(ctx context.Context, filters *models.SaleFilters) ([]*models.Sale, error) {
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
	res, err := s.repo.GetSalesByFilter(ctx, filters)
	if err != nil {
		return nil, errorx.Wrap("get sales", errorx.Internal, err)
	}
	return res, nil
}
