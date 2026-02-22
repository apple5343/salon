package simulation

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

type EventType string

const (
	BrandCreated EventType = "brand_created"

	ModelCreated EventType = "model_created"

	CarCreated EventType = "car_created"
	CarUpdated EventType = "car_updated"

	SaleCreated   EventType = "sale_created"
	SaleCompleted EventType = "sale_completed"
	SaleCancled   EventType = "sale_cancled"

	EmployeeCreated EventType = "employee_created"
	EmployeeUpdated EventType = "employee_updated"
	EmployeeHired   EventType = "employee_hired"
)

const timeLayout = "2006-01-02"

type Event struct {
	Type EventType
	Data interface{}
	Date string
	Time *time.Time
}

type EventNode struct {
	e             *Event
	next          *EventNode
	executionTime time.Duration
}

func (s *Simulation) AddEventNode(node *EventNode) {
	s.eventNodeQueue[node.e.Date] = append(s.eventNodeQueue[node.e.Date], node)
}

func (s *Simulation) AddEvent(e *Event) {
	s.eventQueue[e.Date] = append(s.eventQueue[e.Date], e)
}

func (s *Simulation) handleEventNode(head *EventNode) {
	startHour := 8 + rand.IntN(2)
	startMinute := rand.IntN(60)
	t := s.currendDay.Add(time.Duration(startHour)*time.Hour + time.Duration(startMinute)*time.Minute)
	s.handleEvent(head.e, t)
	node := head.next
	for node != nil && node.executionTime > 0 {
		t = t.Add(node.executionTime)
		s.handleEvent(node.e, t)
		node = node.next
	}
	if node != nil {
		s.AddEventNode(node)
	}
}

func (s *Simulation) handleEvent(e *Event, t time.Time) {
	fmt.Println(e.Type, t)
	s.clock.Set(t)
	switch e.Type {
	case EmployeeCreated:
		s.ProcessEmployeeCreatedEvent(e, t)
	case EmployeeHired:
		s.ProcessHireEmployeeEvent(e, t)
	case BrandCreated:
		s.ProcessBrandCreatedEvent(e, t)
	case ModelCreated:
		s.ProcessModelCreatedEvent(e, t)
	case CarCreated:
		s.ProcessCarCreatedEvent(e, t)
	case CarUpdated:
		s.ProcessCarUpdateEvent(e, t)
	case SaleCreated:
		s.ProcessSaleCreatedEvent(e, t)
	case SaleCompleted:
		s.ProcessSaleCompletedEvent(e, t)
	case SaleCancled:
		s.ProcessSaleCancledEvent(e, t)
	default:
		log.Println("invalid event type")
	}
}
