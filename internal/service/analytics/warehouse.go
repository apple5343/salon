package analytics

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *analyticsService) Warehouse(ctx context.Context, turnoverDateFrom, turnoverDateTo *time.Time) (*models.WarehouseAnalytics, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if turnoverDateFrom != nil && turnoverDateTo != nil && turnoverDateFrom.After(*turnoverDateTo) {
		return nil, ErrInvalidTimeRange
	}
	result, err := s.analyticsRepository.GetWarehouseAnalytics(ctx, turnoverDateFrom, turnoverDateTo)
	if err != nil {
		return nil, errorx.Wrap("get warehouse statistic", errorx.Internal, err)
	}
	return result, nil
}
