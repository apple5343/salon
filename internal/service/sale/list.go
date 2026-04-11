package sale

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
)

func (s *saleService) GetSales(ctx context.Context, filters *models.SaleFilters) ([]*models.Sale, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
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
		orderBy := models.SaleOrderByDate
		filters.OrderBy = &orderBy
	}
	res, err := s.repo.GetSalesByFilter(ctx, filters)
	if err != nil {
		return nil, errorx.Wrap("get sales", errorx.Internal, err)
	}
	return res, nil
}
