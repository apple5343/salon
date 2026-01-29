package car

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden     = errorx.NewError("forbidden", errorx.Forbidden)
	ErrBrandNotFound = errorx.NewError("brand not found", errorx.BadRequest)
	ErrModelNotFound = errorx.NewError("model not found", errorx.BadRequest)
	ErrBrandExists   = errorx.NewError("brand already exists", errorx.Conflict)
	ErrModelExists   = errorx.NewError("model already exists", errorx.Conflict)
	ErrCarExists     = errorx.NewError("car already exists", errorx.Conflict)
	ErrCarNotFound   = errorx.NewError("car not found", errorx.BadRequest)
	ErrForeignKey    = errorx.NewError("foreign key, entity not found", errorx.BadRequest)
)

type carService struct {
	repo         repository.CarRepository
	eventService service.EventService
	supplier     service.SupplierService
}

func NewService(repo repository.CarRepository, eventService service.EventService, supplierService service.SupplierService) service.CarService {
	return &carService{
		repo:         repo,
		eventService: eventService,
		supplier:     supplierService,
	}
}
