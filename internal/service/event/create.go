package event

import (
	"context"
	"log"
	"salon/internal/models"
	"time"
)

const EventProcessingTimeout = 10

func (s *eventService) AddEvent(ctx context.Context, e *models.Event) error {
	processingCtx, cancel := context.WithTimeout(context.Background(), EventProcessingTimeout*time.Second)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer cancel()
		select {
		case <-processingCtx.Done():
			//TODO обработать случай
			log.Println("timeout")
			return
		default:
			err := s.handleEvent(processingCtx, e)
			if nil == err {
				return
			}
			log.Println("event handled", err.Error())
			//TODO обработать случай
		}
	}()
	return nil
}

func (s *eventService) handleEvent(ctx context.Context, e *models.Event) error {
	e, err := s.repo.Create(ctx, e)
	if err != nil {
		return err
	}
	return nil
}
