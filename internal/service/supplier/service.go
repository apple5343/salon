package supplier

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrSupplierNotFound = errorx.NewError("supplier not found", errorx.BadRequest)
	ErrForbidden        = errorx.NewError("forbidden", errorx.Forbidden)
	ErrInvalidID        = errorx.NewError("invalid id", errorx.BadRequest)
)

type supplierService struct {
	repo         repository.SupplierRepository
	cache        repository.SupplierCache
	eventService service.EventService
	clock        clock.Clock
}

func NewService(repo repository.SupplierRepository, eventService service.EventService,
	cache repository.SupplierCache, clock clock.Clock) service.SupplierService {
	return &supplierService{
		repo:         repo,
		cache:        cache,
		eventService: eventService,
		clock:        clock,
	}
}
