package car

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden   = errorx.NewError("forbidden", errorx.Forbidden)
	ErrCarExists   = errorx.NewError("car already exists", errorx.Conflict)
	ErrCarNotFound = errorx.NewError("car not found", errorx.BadRequest)
	ErrForeignKey  = errorx.NewError("foreign key, entity not found", errorx.BadRequest)
	ErrInvalidID   = errorx.NewError("invalid id", errorx.BadRequest)
)

type carService struct {
	repo         repository.CarRepository
	eventService service.EventService
	supplier     service.SupplierService
	model        service.ModelService
	brand        service.BrandService
	clock        clock.Clock
}

func NewService(repo repository.CarRepository, eventService service.EventService,
	supplierService service.SupplierService, model service.ModelService,
	brand service.BrandService, clock clock.Clock) service.CarService {
	return &carService{
		repo:         repo,
		eventService: eventService,
		supplier:     supplierService,
		model:        model,
		brand:        brand,
		clock:        clock,
	}
}
