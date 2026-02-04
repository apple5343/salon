package sale

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden    = errorx.NewError("forbidden", errorx.Forbidden)
	ErrSaleNotFound = errorx.NewError("sale not found", errorx.BadRequest)
)

type saleService struct {
	repo       repository.SaleRepository
	carService service.CarService
}

func NewService(repo repository.SaleRepository, carService service.CarService) service.SaleService {
	return &saleService{
		repo: repo,
		carService: carService,
	}
}
