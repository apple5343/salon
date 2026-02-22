package client

import (
	"salon/internal/repository"
	"salon/internal/service"
	"salon/pkg/clock"

	"github.com/apple5343/errorx"
)

var (
	ErrUnauthorized   = errorx.NewError("unauthorized", errorx.Unauthorized)
	ErrForbidden      = errorx.NewError("forbidden", errorx.Forbidden)
	ErrClientExists   = errorx.NewError("client already exists", errorx.Conflict)
	ErrClientNotFound = errorx.NewError("client not found", errorx.BadRequest)
	ErrInvalidID      = errorx.NewError("invalid id", errorx.BadRequest)
)

type clientService struct {
	repo  repository.ClientRepository
	clock clock.Clock
}

func NewService(repo repository.ClientRepository, clock clock.Clock) service.ClientService {
	return &clientService{
		repo:  repo,
		clock: clock,
	}
}
