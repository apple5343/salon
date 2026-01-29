package http

import (
	"context"
	"salon/internal/config"
	"salon/internal/service"
	"salon/internal/transport/http/car"
	"salon/internal/transport/http/client"
	"salon/internal/transport/http/employee"
	"salon/internal/transport/http/event"
	"salon/internal/transport/http/middlewares"
	"salon/internal/transport/http/supplier"
	"salon/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	config          *config.HttpServer
	jwtConfig       *config.JWT
	logger          logger.Logger
	employeeHandler *employee.Handler
	supplierHandler *supplier.Handler
	clientHandler   *client.Handler
	carHandler      *car.Handler
	eventHandler    *event.Handler
	e               *echo.Echo
}

func NewServer(config *config.HttpServer, jwtConfig *config.JWT, logger logger.Logger,
	employeeService service.EmployeeService, supplierService service.SupplierService,
	clientService service.ClientService, carService service.CarService,
	eventService service.EventService) (*Server, error) {
	return &Server{
		config:          config,
		jwtConfig:       jwtConfig,
		employeeHandler: employee.NewHandler(employeeService),
		supplierHandler: supplier.NewHandler(supplierService),
		clientHandler:   client.NewHandler(clientService),
		carHandler:      car.NewHandler(carService),
		eventHandler:    event.NewHandler(eventService),
		logger:          logger,
		e:               echo.New(),
	}, nil
}

func (s *Server) Start(_ context.Context) error {
	s.routes()
	return s.e.Start(s.config.Addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) routes() {
	authMiddleware := middlewares.AuthMiddleware(s.jwtConfig)
	softAuthMiddleware := middlewares.SoftAuthMiddleware(s.jwtConfig)
	tokenMiddleware := middlewares.TokenMiddleware()
	s.e.Use(middlewares.LoggerMiddleware(s.logger))
	s.e.Use(middlewares.ErrorMiddleware())

	employees := s.e.Group("/employees")
	{
		auth := employees.Group("/auth")
		{
			auth.POST("/login", s.employeeHandler.Login())
			auth.POST("/register", s.employeeHandler.Register()) //TODO добавить мидлвары
			auth.GET("/access", tokenMiddleware(s.employeeHandler.GetAccessToken()))
			auth.GET("/refresh", tokenMiddleware(s.employeeHandler.GetRefreshToken()))
		}
	}

	clients := s.e.Group("clients")
	{
		auth := clients.Group("/auth")
		{
			auth.POST("/register", authMiddleware(s.clientHandler.Register()))
		}
		clients.GET("/:id", authMiddleware(s.clientHandler.GetByID()))
		clients.PUT("/:id", authMiddleware(s.clientHandler.Update()))
	}

	brands := s.e.Group("/brands")
	{
		brands.POST("", authMiddleware(s.carHandler.CreateBrand()))
		brands.GET("/:id", softAuthMiddleware(s.carHandler.GetBrandByID()))
		brands.PUT("/:id", authMiddleware(s.carHandler.UpdateBrand()))
	}

	models := s.e.Group("/models")
	{
		models.POST("", authMiddleware(s.carHandler.CreateModel()))
		models.GET("", s.carHandler.GetModels())
		models.GET("/:id", softAuthMiddleware(s.carHandler.GetModelByID()))
		models.PUT("/:id", authMiddleware(s.carHandler.UpdateModel()))
	}

	cars := s.e.Group("/cars")
	{
		cars.POST("", authMiddleware(s.carHandler.CreateCar()))
		cars.GET("", s.carHandler.GetCars())
		cars.GET("/:id", softAuthMiddleware(s.carHandler.GetCarByID()))
		cars.PUT("/:id", authMiddleware(s.carHandler.UpdateCar()))
	}

	suppliers := s.e.Group("/suppliers")
	{
		suppliers.POST("", authMiddleware(s.supplierHandler.Create()))
		suppliers.GET("/:id", softAuthMiddleware(s.supplierHandler.GetByID()))
		suppliers.PUT("/:id", authMiddleware(s.supplierHandler.Update()))
	}

	events := s.e.Group("/events")
	{
		events.GET("", authMiddleware(s.eventHandler.GetEvents()))
		events.GET("/:id", authMiddleware(s.eventHandler.GetByID()))
	}
}
