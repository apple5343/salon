package brand

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrInvalidID     = errorx.NewError("invalid id", errorx.BadRequest)
	ErrBrandNotFound = errorx.NewError("brand not found", errorx.BadRequest)
	ErrForbidden     = errorx.NewError("forbidden", errorx.Forbidden)
	ErrBrandExists   = errorx.NewError("brand already exists", errorx.Conflict)
)

type brandService struct {
	repo         repository.BrandRepository
	eventService service.EventService
	clock        clock.Clock
}

func NewService(repo repository.BrandRepository, eventService service.EventService, clock clock.Clock) service.BrandService {
	return &brandService{
		repo:         repo,
		eventService: eventService,
		clock:        clock,
	}
}
