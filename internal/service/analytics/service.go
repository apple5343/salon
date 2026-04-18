package analytics

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden = errorx.NewError("forbidden", errorx.Forbidden)
)

type analyticsService struct {
	analyticsRepository repository.AnalyticsRepository
	employeeService     service.EmployeeService
}

func NewService(repo repository.AnalyticsRepository, employeeService service.EmployeeService) service.AnalyticsService {
	return &analyticsService{
		analyticsRepository: repo,
		employeeService:     employeeService,
	}
}
