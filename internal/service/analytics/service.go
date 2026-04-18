package analytics

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden        = errorx.NewError("forbidden", errorx.Forbidden)
	ErrInvalidTimeRange = errorx.NewError("invalid time range", errorx.BadRequest)
)

type analyticsService struct {
	analyticsRepository repository.AnalyticsRepository
	employeeService     service.EmployeeService
	supplierService     service.SupplierService
}

func NewService(repo repository.AnalyticsRepository, employeeService service.EmployeeService,
	supplierService service.SupplierService) service.AnalyticsService {
	return &analyticsService{
		analyticsRepository: repo,
		employeeService:     employeeService,
		supplierService:     supplierService,
	}
}
