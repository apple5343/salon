package event

import (
	"context"
	"salon/internal/repository"
	"salon/internal/service"
	"sync"

	"github.com/apple5343/errorx"
)

var (
	ErrEventNotFound = errorx.NewError("event not found", errorx.BadRequest)
	ErrForbidden     = errorx.NewError("forbidden", errorx.Forbidden)
)

type eventService struct {
	repo repository.EventRepository
	wg   sync.WaitGroup
}

func NewService(repo repository.EventRepository) service.EventService {
	return &eventService{
		repo: repo,
	}
}

func (s *eventService) Shutdown(_ context.Context) error {
	s.wg.Wait()
	return nil
}
