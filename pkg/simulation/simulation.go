package simulation

import (
	"math/rand/v2"
	"salon/internal/models"
	"salon/internal/service"
	"salon/pkg/clock"
	"salon/pkg/simulation/generator"
	"sort"
	"time"
)

type Simulation struct {
	cfg              *Config
	currendDay       time.Time
	activeEmployees  []string
	activeAdmins     []string
	availableClients []string
	availableCars    map[string]bool
	eventNodeQueue   map[string][]*EventNode
	eventQueue       map[string][]*Event
	employees        map[string]*models.Employee
	admins           map[string]*models.Employee
	cars             map[string]*models.Car

	carService      service.CarService
	employeeService service.EmployeeService

	generator *generator.Generator
	clock     clock.MockClock
}

func NewSimulation(carService service.CarService, employeeService service.EmployeeService, clock clock.MockClock, cfg *Config) *Simulation {
	return &Simulation{
		cfg:             cfg,
		eventNodeQueue:  make(map[string][]*EventNode),
		eventQueue:      make(map[string][]*Event),
		employees:       make(map[string]*models.Employee),
		admins:          make(map[string]*models.Employee),
		cars:            make(map[string]*models.Car),
		availableCars:   make(map[string]bool),
		carService:      carService,
		employeeService: employeeService,
		clock:           clock,
	}
}

func (s *Simulation) Run() {
	s.Init()
	for i := 0; i < s.cfg.DaysCount; i++ {
		s.PlanDay()
		s.ProcessDay()
		s.currendDay = s.currendDay.AddDate(0, 0, 1)
	}
}

func (s *Simulation) ProcessDay() {
	type event struct {
		e    interface{}
		time time.Time
	}
	dateKey := s.currendDay.Format(timeLayout)
	q := []event{}
	for _, node := range s.eventNodeQueue[dateKey] {
		q = append(q, event{node, s.RandomDayTime()})
	}
	for _, e := range s.eventQueue[dateKey] {
		q = append(q, event{e, s.RandomDayTime()})
	}
	sort.Slice(q, func(i, j int) bool {
		return q[i].time.Before(q[j].time)
	})

	for _, e := range q {
		switch e.e.(type) {
		case *Event:
			s.handleEvent(e.e.(*Event), e.time)
		case *EventNode:
			s.handleEventNode(e.e.(*EventNode))
		}
	}
}

func (s *Simulation) Init() {
	s.currendDay = s.cfg.StartDate
	for i := 0; i < 3; i++ {
		s.CreateAdmin(s.currendDay)
	}
	for i := 0; i < 5; i++ {
		s.HireEmployeeRandom()
	}
	s.ProcessDay()
	s.currendDay = s.cfg.StartDate
}

func (s *Simulation) PlanDay() {
	if s.currendDay.Weekday() == time.Saturday || s.currendDay.Weekday() == time.Sunday {
		return
	}
	newEmployees := Poisson(0.1)
	for i := 0; i < newEmployees; i++ {
		s.HireEmployeeRandom()
	}
}

func (s *Simulation) RandomDayTime() time.Time {
	hour := 8 + rand.IntN(12)
	minute := rand.IntN(60)
	return s.currendDay.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
}
