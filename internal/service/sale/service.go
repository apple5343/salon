package sale

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden       = errorx.NewError("forbidden", errorx.Forbidden)
	ErrSaleNotFound    = errorx.NewError("sale not found", errorx.BadRequest)
	ErrInvalidID       = errorx.NewError("invalid id", errorx.BadRequest)
	ErrForeignKey      = errorx.NewError("no entity", errorx.BadRequest)
	ErrCarNotFound     = errorx.NewError("car not found", errorx.BadRequest)
	ErrCarNotAvailable = errorx.NewError("car is not available", errorx.BadRequest)
)

type saleService struct {
	repo          repository.SaleRepository
	carService    service.CarService
	clientService service.ClientService
	clock         clock.Clock
}

func NewService(repo repository.SaleRepository, carService service.CarService,
	clientService service.ClientService, clock clock.Clock) service.SaleService {
	return &saleService{
		repo:          repo,
		carService:    carService,
		clientService: clientService,
		clock:         clock,
	}
}
