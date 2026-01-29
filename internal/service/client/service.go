package client

import (
	"salon/internal/repository"
	"salon/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrUnauthorized   = errorx.NewError("unauthorized", errorx.Unauthorized)
	ErrForbidden      = errorx.NewError("forbidden", errorx.Forbidden)
	ErrClientExists   = errorx.NewError("client already exists", errorx.Conflict)
	ErrClientNotFound = errorx.NewError("client not found", errorx.BadRequest)
)

type clientService struct {
	repo repository.ClientRepository
}

func NewService(repo repository.ClientRepository) service.ClientService {
	return &clientService{
		repo: repo,
	}
}
