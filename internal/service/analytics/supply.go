package analytics

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *analyticsService) Supply(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SupplyAnalytics, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if dateFrom != nil && dateTo != nil && dateFrom.After(*dateTo) {
		return nil, ErrInvalidTimeRange
	}
	result, err := s.analyticsRepository.GetSupplyAnalytics(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, errorx.Wrap("get supply statistic", errorx.Internal, err)
	}
	return result, nil
}
