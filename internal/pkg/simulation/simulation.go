package simulation

import (
	"errors"
	"math/rand/v2"
	"salon/internal/models"
	"salon/internal/pkg/simulation/generator"
	"salon/internal/service"
	"salon/pkg/clock"
	"sort"
	"time"
)

var (
	ErrNoAvailableAdmins    = errors.New("no available admins")
	ErrNoAvailableEmployees = errors.New("no available employees")
	ErrNoAvailableCars      = errors.New("no available cars")
)

type Simulation struct {
	cfg        *Config
	currendDay time.Time

	activeEmployees  []string
	activeAdmins     []string
	availableClients []string

	availableCars  map[string]bool
	eventNodeQueue map[string][]*EventNode
	eventQueue     map[string][]*Event

	employees map[string]*models.Employee
	clients   map[string]*models.Client
	admins    map[string]*models.Employee
	sales     map[string]*models.Sale

	modelService    service.ModelService
	brandService    service.BrandService
	carService      service.CarService
	employeeService service.EmployeeService
	clientService   service.ClientService
	saleService     service.SaleService
	supplierService service.SupplierService

	generator *generator.Generator
	clock     clock.MockClock
}

func NewSimulation(carService service.CarService, modelService service.ModelService,
	brandService service.BrandService, employeeService service.EmployeeService,
	supplierService service.SupplierService, clientService service.ClientService,
	saleService service.SaleService,
	clock clock.MockClock, cfg *Config) *Simulation {
	return &Simulation{
		cfg:             cfg,
		eventNodeQueue:  make(map[string][]*EventNode),
		eventQueue:      make(map[string][]*Event),
		employees:       make(map[string]*models.Employee),
		admins:          make(map[string]*models.Employee),
		clients:         make(map[string]*models.Client),
		sales:           make(map[string]*models.Sale),
		availableCars:   make(map[string]bool),
		carService:      carService,
		employeeService: employeeService,
		brandService:    brandService,
		modelService:    modelService,
		clientService:   clientService,
		supplierService: supplierService,
		saleService:     saleService,
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
	println("Employees total: ", len(s.employees), " Admins total: ", len(s.admins), " Cars total: ", s.generator.CreatedCarsCount(),
		" Brands total: ", s.generator.CreatedBrandsCount(), " Models total: ", s.generator.CreatedModelsCount(),
		" Suppliers total: ", s.generator.CreatedSuppliersCount(), " Clients total: ", len(s.clients), " Sales total: ", len(s.sales))
}

func (s *Simulation) Init() {
	s.currendDay = s.cfg.StartDate
	for i := 0; i < s.cfg.AdminsCount; i++ {
		s.CreateAdmin(s.currendDay)
	}
	for i := 0; i < s.cfg.EmployeesCount; i++ {
		s.CreateEmployee()
	}
	suppliersCount := s.generator.AvailableSuppliersCount() * s.cfg.SuppliersPercent / 100
	for i := 0; i < suppliersCount; i++ {
		s.CreateSupplier()
	}
	brandsCount := s.generator.AvailableBrandsCount() * s.cfg.BrandsPercent / 100
	for i := 0; i < brandsCount; i++ {
		s.CreateBrand()
	}
	modelsCount := s.generator.PendingModelsCount() * s.cfg.ModelsPercent / 100
	for i := 0; i < modelsCount; i++ {
		s.CreateModel()
	}
	s.ProcessDay()
	s.currendDay = s.currendDay.AddDate(0, 0, 1)
}

func (s *Simulation) PlanDay() {
	if s.currendDay.Weekday() == time.Saturday || s.currendDay.Weekday() == time.Sunday {
		return
	}
	newEmployees := Poisson(s.cfg.NewEmloyeesRatio)
	for i := 0; i < newEmployees; i++ {
		s.CreateEmployee()
	}
	newCars := Poisson(s.cfg.NewCarsRatio)
	for i := 0; i < newCars; i++ {
		s.CreateCar()
	}
	newBrands := Poisson(s.cfg.NewBrandsRatio)
	for i := 0; i < newBrands; i++ {
		s.CreateBrand()
	}
	newModels := Poisson(s.cfg.NewModelsRatio)
	for i := 0; i < newModels; i++ {
		s.CreateModel()
	}
	newSuppliers := Poisson(s.cfg.NewSuppliersRatio)
	for i := 0; i < newSuppliers; i++ {
		s.CreateSupplier()
	}
	newSales := Poisson(s.cfg.NewSalesRatio)
	for i := 0; i < newSales; i++ {
		s.CreateSale()
	}
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

func (s *Simulation) RandomDayTime() time.Time {
	hour := 8 + rand.IntN(12)
	minute := rand.IntN(60)
	seconds := rand.IntN(60)
	milliseconds := rand.IntN(1000)
	return s.currendDay.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond)
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
