package analytics

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *analyticsService) Supplier(ctx context.Context, supplierID string, dateFrom, dateTo *time.Time) (*models.SupplierAnalytics, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if _, err := s.supplierService.GetByID(ctx, supplierID); err != nil {
		return nil, err
	}

	if dateFrom != nil && dateTo != nil && dateFrom.After(*dateTo) {
		return nil, ErrInvalidTimeRange
	}
	result, err := s.analyticsRepository.GetSupplierAnalytics(ctx, supplierID, dateFrom, dateTo)
	if err != nil {
		return nil, errorx.Wrap("get supplier statistic", errorx.Internal, err)
	}
	return result, nil
}
