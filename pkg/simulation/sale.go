package simulation

import (
	"context"
	"log"
	"math/rand/v2"
	ctxutil "salon/internal/utils/context"
	"time"
)

type SaleCreatedEvent struct {
	CarID       string
	ClientID    string
	EmployeeID  string
	IsCompleted bool
}

type SaleCompletedEvent struct {
	EmployeeID string
	SaleID     string
}

type SaleCanceledEvent struct {
	EmployeeID string
	SaleID     string
}

func (s *Simulation) CreateSale() {
	s.AddEvent(&Event{
		Type: ClientCreated,
		Date: s.currendDay.Format(timeLayout),
	})
}

func (s *Simulation) ProcessClientCreatedEvent(e *Event, t time.Time) {
	employeeID, ok := s.RandomEmployee()
	if !ok {
		log.Println("create sale: " + ErrNoAvailableEmployees.Error())
		return
	}
	car, ok := s.RandomAvailableCar()
	if !ok {
		log.Println("create sale: " + ErrNoAvailableCars.Error())
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), employeeID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(s.employees[employeeID].Role))
	client, err := s.clientService.Register(ctx, s.generator.GenerateClient())
	if err != nil {
		log.Println("register client: " + err.Error())
		return
	}
	s.clients[client.ID] = client
	s.availableCars[car.ID] = false
	saleCreatedTime := t.Add(s.RandomDurationMinutes(15, 120))
	isCompleted := true
	if rand.IntN(10) < 3 {
		isCompleted = false
	}
	s.AddEvent(&Event{
		Type: SaleCreated,
		Data: SaleCreatedEvent{
			CarID:       car.ID,
			ClientID:    client.ID,
			EmployeeID:  employeeID,
			IsCompleted: isCompleted,
		},
		Date: t.Format(timeLayout),
		Time: &saleCreatedTime,
	})
}

func (s *Simulation) ProcessSaleCreatedEvent(e *Event, t time.Time) {
	data, ok := e.Data.(SaleCreatedEvent)
	if !ok {
		log.Println("invalid event data")
		return
	}
	employee, ok := s.employees[data.EmployeeID]
	if !ok {
		log.Println("create sale: employee not found")
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), employee.ID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(employee.Role))
	car, ok := s.generator.GetCreatedCar(data.CarID)
	if !ok {
		log.Println("create sale: car not found")
		return
	}
	client, ok := s.clients[data.ClientID]
	if !ok {
		log.Println("create sale: client not found")
		return
	}
	sale, err := s.saleService.Create(ctx, s.generator.GenerateSale(car, client.ID, employee.ID))
	if err != nil {
		log.Println("create sale: " + err.Error())
		return
	}
	s.sales[sale.ID] = sale
	nextTime := t.Add(s.RandomDurationMinutes(15, 120))
	if data.IsCompleted {
		s.AddEvent(&Event{
			Type: SaleCompleted,
			Data: SaleCompletedEvent{
				EmployeeID: data.EmployeeID,
				SaleID:     sale.ID,
			},
			Date: t.Format(timeLayout),
			Time: &nextTime,
		})
	} else {
		s.AddEvent(&Event{
			Type: SaleCanceled,
			Data: SaleCanceledEvent{
				EmployeeID: data.EmployeeID,
				SaleID:     sale.ID,
			},
			Date: t.Format(timeLayout),
			Time: &nextTime,
		})
	}
}

func (s *Simulation) ProcessSaleCanceledEvent(e *Event, t time.Time) {
	data, ok := e.Data.(SaleCanceledEvent)
	if !ok {
		log.Println("invalid event data")
		return
	}
	employee, ok := s.employees[data.EmployeeID]
	if !ok {
		log.Println("cancel sale: employee not found")
		return
	}
	sale, ok := s.sales[data.SaleID]
	if !ok {
		log.Println("cancel sale: sale not found")
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), employee.ID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(employee.Role))
	if err := s.saleService.Cancel(ctx, data.SaleID); err != nil {
		log.Println("cancel sale: " + err.Error())
		return
	}
	sale, err := s.saleService.GetByID(ctx, data.SaleID)
	if err != nil {
		log.Println("cancel sale: " + err.Error())
		return
	}
	s.availableCars[sale.CarID] = true
	s.sales[data.SaleID] = sale
}

func (s *Simulation) ProcessSaleCompletedEvent(e *Event, t time.Time) {
	data, ok := e.Data.(SaleCompletedEvent)
	if !ok {
		log.Println("invalid event data")
		return
	}
	employee, ok := s.employees[data.EmployeeID]
	if !ok {
		log.Println("complete sale: employee not found")
		return
	}
	_, ok = s.sales[data.SaleID]
	if !ok {
		log.Println("complete sale: sale not found")
		return
	}
	ctx := ctxutil.ContextWithUserID(context.TODO(), employee.ID)
	ctx = ctxutil.ContextWithUserRole(ctx, string(employee.Role))
	sale, err := s.saleService.Complete(ctx, data.SaleID)
	if err != nil {
		log.Println("complete sale: " + err.Error())
		return
	}
	s.sales[data.SaleID] = sale
}
