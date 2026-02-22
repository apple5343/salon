package model

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrForbidden     = errorx.NewError("forbidden", errorx.Forbidden)
	ErrModelNotFound = errorx.NewError("model not found", errorx.BadRequest)
	ErrModelExists   = errorx.NewError("model already exists", errorx.Conflict)
	ErrForeignKey    = errorx.NewError("foreign key, entity not found", errorx.BadRequest)
	ErrInvalidID     = errorx.NewError("invalid id", errorx.BadRequest)
)

type modelService struct {
	repo         repository.ModelRepository
	brandService service.BrandService
	eventService service.EventService
	clock        clock.Clock
}

func NewService(repo repository.ModelRepository, eventService service.EventService,
	brandService service.BrandService, clock clock.Clock) service.ModelService {
	return &modelService{
		repo:         repo,
		brandService: brandService,
		eventService: eventService,
		clock:        clock,
	}
}
