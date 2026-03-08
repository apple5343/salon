package simulation

import (
	"errors"
	"math/rand/v2"
	"salon/internal/models"
	"salon/internal/service"
	"salon/pkg/clock"
	"salon/pkg/simulation/generator"
	"sort"
	"time"
)

var (
	ErrNoAvailableAdmins = errors.New("no available admins")
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

	modelService    service.ModelService
	brandService    service.BrandService
	carService      service.CarService
	employeeService service.EmployeeService

	generator *generator.Generator
	clock     clock.MockClock
}

func NewSimulation(carService service.CarService, modelService service.ModelService, brandService service.BrandService, employeeService service.EmployeeService, clock clock.MockClock, cfg *Config) *Simulation {
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
		brandService:    brandService,
		modelService:    modelService,
		generator:       generator.NewGenerator(),
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
	println("Employees total: ", len(s.employees), " Admins total: ", len(s.admins), " Cars total: ", len(s.cars), " Brands total: ", s.generator.CreatedBrandsCount(), " Models total: ", s.generator.CreatedModelsCount())
}

func (s *Simulation) ProcessDay() {
	type event struct {
		e    interface{}
		time time.Time
	}
	dateKey := s.currendDay.Format(timeLayout)
	nodeQueue := s.eventNodeQueue[dateKey]
	eventQueue := s.eventQueue[dateKey]
	if len(nodeQueue) != 0 || len(eventQueue) != 0 {
		println(dateKey)
	}
	s.eventNodeQueue[dateKey] = nil
	s.eventQueue[dateKey] = nil
	for len(nodeQueue) > 0 || len(eventQueue) > 0 {
		q := []event{}
		for _, node := range nodeQueue {
			q = append(q, event{node, s.RandomDayTime()})
		}
		for _, e := range eventQueue {
			if e.Time != nil {
				q = append(q, event{e, *e.Time})
				continue
			}
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
		nodeQueue = s.eventNodeQueue[dateKey]
		eventQueue = s.eventQueue[dateKey]
		s.eventNodeQueue[dateKey] = nil
		s.eventQueue[dateKey] = nil
	}
	delete(s.eventNodeQueue, dateKey)
	delete(s.eventQueue, dateKey)
}

func (s *Simulation) Init() {
	s.currendDay = s.cfg.StartDate
	for i := 0; i < 3; i++ {
		s.CreateAdmin(s.currendDay)
	}
	for i := 0; i < 5; i++ {
		s.CreateEmployee()
	}
	s.ProcessDay()
	s.currendDay = s.currendDay.AddDate(0, 0, 1)
}

func (s *Simulation) PlanDay() {
	if s.currendDay.Weekday() == time.Saturday || s.currendDay.Weekday() == time.Sunday {
		return
	}
	newEmployees := Poisson(0.01)
	for i := 0; i < newEmployees; i++ {
		s.CreateEmployee()
	}
	newCars := Poisson(0.01)
	for i := 0; i < newCars; i++ {
		s.CreateCar()
	}
	newBrands := Poisson(0.07)
	for i := 0; i < newBrands; i++ {
		s.CreateBrand()
	}
	newModels := Poisson(0.07)
	for i := 0; i < newModels; i++ {
		s.CreateModel()
	}
}

func (s *Simulation) RandomDayTime() time.Time {
	hour := 8 + rand.IntN(12)
	minute := rand.IntN(60)
	return s.currendDay.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
}

func (s *Simulation) RandomDurationMinutes(minMinutes int, maxMinutes int) time.Duration {
	minutes := minMinutes + rand.IntN(maxMinutes-minMinutes+1)
	seconds := rand.IntN(60)
	milliseconds := rand.IntN(1000)
	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond
}

func (s *Simulation) RandomInt(min int, max int) int {
	return min + rand.IntN(max-min+1)
}
