package simulation

import (
	"context"
	"log"
	"math/rand/v2"
	"salon/internal/models"
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
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	car := s.generator.GenerateCar()
	// car, _, _, _, err := s.carService.CreateCar(ctx, car)
	// if err != nil {
	// 	log.Println("create car: " + err.Error())
	// 	return
	// }
	s.cars[car.ID] = &car
	s.GenerateCarAvailableFlow(&car, t.AddDate(0, 0, s.RandomInt(7, 20)))
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
	// ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	// ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	// car, ok := s.cars[data.CarID]
	// if !ok {
	// 	log.Println("update car: car not found")
	// 	return
	// }
	// car, err := s.carService.UpdateCar(ctx)
	// s.carService.UpdateCar()
}

func (s *Simulation) RandomAvailableCar() *models.Car {
	if len(s.availableCars) == 0 {
		return nil
	}
	keys := make([]string, 0, len(s.availableCars))
	for k := range s.availableCars {
		keys = append(keys, k)
	}
	return s.cars[keys[rand.IntN(len(keys))]]
}
