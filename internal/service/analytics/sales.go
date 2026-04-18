package analytics

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *analyticsService) Sales(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SalesAnalytics, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	result, err := s.analyticsRepository.GetSalesAnalytics(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, errorx.Wrap("get sales statistic", errorx.Internal, err)
	}
	return result, nil
}
