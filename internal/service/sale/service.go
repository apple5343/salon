package sale

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden    = errorx.NewError("forbidden", errorx.Forbidden)
	ErrSaleNotFound = errorx.NewError("sale not found", errorx.BadRequest)
)

type saleService struct {
	repo       repository.SaleRepository
	carService service.CarService
	clock      clock.Clock
}

func NewService(repo repository.SaleRepository, carService service.CarService, clock clock.Clock) service.SaleService {
	return &saleService{
		repo:       repo,
		carService: carService,
		clock:      clock,
	}
}
