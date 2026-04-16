package simulation

import (
	"context"
	"errors"
	"log"
	"salon/internal/models"
	"salon/internal/pkg/simulation/generator"
	ctxutil "salon/internal/utils/context"
	"time"

	"github.com/brianvoe/gofakeit"
)

type CarUpdateStatusEvent struct {
	CarID  string
	Status models.CarStatus
}

func (s *Simulation) CreateCar() {
	s.AddEvent(&Event{
		Type: CarCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessCarCreatedEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("create car: " + ErrNoAvailableAdmins.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	car, err := s.generator.GenerateCar()
	if err != nil {
		if errors.Is(err, generator.ErrNoItems) {
			return
		}
		log.Println("create car: " + err.Error())
		return
	}
	created, _, _, _, err := s.carService.CreateCar(ctx, &car)
	if err != nil {
		log.Println("create car with id " + car.ID + ":" + err.Error())
		return
	}
	s.generator.CarCreated(created)
	s.GenerateCarAvailableFlow(created, t.AddDate(0, 0, s.RandomInt(7, 20)))
}

func (s *Simulation) GenerateCarAvailableFlow(car *models.Car, startTime time.Time) {
	date := startTime
	var head, node *EventNode
	addNode := func(event *Event, executionTime time.Duration) {
		n := &EventNode{e: event, executionTime: executionTime}
		if head == nil {
			head = n
			node = head
		} else {
			node.next = n
			node = node.next
		}
	}
	if car.Status == models.CarStatusIncoming {
		date = date.AddDate(0, 0, gofakeit.Number(3, 10))
		event := &Event{
			Type: CarUpdated,
			Data: CarUpdateStatusEvent{
				CarID:  car.ID,
				Status: models.CarStatusPending,
			},
			Date: date.Format(timeLayout),
		}
		addNode(event, 0)
	}
	date = date.AddDate(0, 0, gofakeit.Number(1, 7))
	event := &Event{
		Type: CarUpdated,
		Data: CarUpdateStatusEvent{
			CarID:  car.ID,
			Status: models.CarStatusAvailable,
		},
		Date: date.Format(timeLayout),
	}
	addNode(event, 0)
	s.AddEventNode(head)
}

func (s *Simulation) ProcessCarUpdateEvent(e *Event, t time.Time) {
	adminID, ok := s.RandomAdmin()
	if !ok {
		log.Println("update car: " + ErrNoAvailableAdmins.Error())
		return
	}
	data, ok := e.Data.(CarUpdateStatusEvent)
	if !ok {
		log.Println("invalid event data")
	}
	_, ok = s.RandomAdmin()
	if !ok {
		log.Println("update car: " + ErrNoAvailableAdmins.Error())
	}
	if data.Status == models.CarStatusAvailable {
		s.availableCars[data.CarID] = true
	}
	car, ok := s.generator.GetCreatedCar(data.CarID)
	if !ok {
		log.Println("update car: car not found")
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	car.Status = data.Status
	car, _, _, _, err := s.carService.UpdateCar(ctx, car)
	if err != nil {
		log.Println("update car: " + err.Error())
		return
	}
	s.generator.CarUpdated(car)
}

func (s *Simulation) RandomAvailableCar() (*models.Car, bool) {
	for id, available := range s.availableCars {
		if available {
			return s.generator.GetCreatedCar(id)
		}
	}
	return nil, false
}
