package simulation

import (
	"log"
	"math/rand/v2"
	"salon/internal/models"
	"time"
)

type CreateSaleEvent struct {
	CarID string
}

type SaleCompletedEvent struct {
	CarID string
}

type SaleCanceledEvent struct {
	CarID string
}

func (s *Simulation) GenerateSaleCompletedFlow(car *models.Car, startDay time.Time) {
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

	event := &Event{
		Type: SaleCreated,
		Data: CreateSaleEvent{
			CarID: car.ID,
		},
		Date: date.Format(timeLayout),
	}
	addNode(event, 0)

	event = &Event{
		Type: SaleCompleted,
		Data: SaleCompletedEvent{
			CarID: car.ID,
		},
		Date: date.Format(timeLayout),
	}
	addNode(event, time.Duration(rand.IntN(30))*time.Minute)
	s.AddEventNode(head)
}

func (s *Simulation) GenerateSaleCanceledFlow(car *models.Car, startDay time.Time) {
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

	event := &Event{
		Type: SaleCreated,
		Data: CreateSaleEvent{
			CarID: car.ID,
		},
		Date: date.Format(timeLayout),
	}
	addNode(event, 0)

	event = &Event{
		Type: SaleCancled,
		Data: SaleCanceledEvent{
			CarID: car.ID,
		},
		Date: date.Format(timeLayout),
	}
	addNode(event, time.Duration(rand.IntN(30))*time.Minute)
	s.AddEventNode(head)
}

func (s *Simulation) ProcessSaleCreatedEvent(e *Event, t time.Time) {
	data, ok := e.Data.(CreateSaleEvent)
	if !ok {
		log.Println("invalid event data")
	}
	delete(s.availableCars, data.CarID)
}

func (s *Simulation) ProcessSaleCancledEvent(e *Event, t time.Time) {
	data, ok := e.Data.(SaleCanceledEvent)
	if !ok {
		log.Println("invalid event data")
	}
	s.availableCars[data.CarID] = true
}

func (s *Simulation) ProcessSaleCompletedEvent(e *Event, t time.Time) {
	_, ok := e.Data.(SaleCompletedEvent)
	if !ok {
		log.Println("invalid event data")
	}
}
