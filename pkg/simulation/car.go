package simulation

import (
	"context"
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

func (s *Simulation) CreateCar(car *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error) {
	adminsIds := make([]string, 0, len(s.admins))
	for id := range s.admins {
		adminsIds = append(adminsIds, id)
	}
	adminID := adminsIds[rand.IntN(len(s.admins))]
	ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	return s.carService.CreateCar(ctx, car)
}

func (s *Simulation) GenerateCarAvailableFlow(car *models.Car, startDay time.Time) {
	date := startDay
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
	// data, ok := e.Data.(CarUpdateStatusEvent)
	// if !ok {
	// 	log.Println("invalid event data")
	// }
	// adminID, ok := s.RandomAdmin()
	// if !ok {
	// 	log.Println("update car: admin not found")
	// 	return
	// }
	// // ctx := ctxutil.ContextWithUserID(context.TODO(), adminID)
	// // ctx = ctxutil.ContextWithUserRole(ctx, string(models.EmployeeRoleAdmin))
	// // car, ok := s.cars[data.CarID]
	// // if !ok {
	// // 	log.Println("update car: car not found")
	// // 	return
	// // }
	// // car, err := s.carService.UpdateCar(ctx)
	// // s.carService.UpdateCar()
	return
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
