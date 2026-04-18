package analytics

import (
	"context"
	"salon/internal/models"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/apple5343/errorx"
)

func (s *analyticsService) Employee(ctx context.Context, employeeID string, dateFrom, dateTo *time.Time) (*models.EmployeeAnalytics, error) {
	if role := ctxutil.UserRoleFromContext(ctx); role != string(models.EmployeeRoleAdmin) {
		return nil, ErrForbidden
	}
	if _, err := s.employeeService.GetByID(ctx, employeeID); err != nil {
		return nil, err
	}
	
	result, err := s.analyticsRepository.GetEmployeeAnalytics(ctx, employeeID, dateFrom, dateTo)
	if err != nil {
		return nil, errorx.Wrap("get employee statistic", errorx.Internal, err)
	}
	return result, nil
}
