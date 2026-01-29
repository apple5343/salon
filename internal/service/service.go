package service

import (
	"context"
	"salon/internal/models"
)

type EmployeeService interface {
	Register(ctx context.Context, e *models.Employee) (*models.Employee, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	GetRefreshToken(ctx context.Context) (string, error)
	GetAccessToken(ctx context.Context) (string, error)
	GetByID(ctx context.Context, id string) (*models.Employee, error)
}

type ClientService interface {
	GetByID(ctx context.Context, id string) (*models.Client, error)
	Register(ctx context.Context, c *models.Client) (*models.Client, error)
	Update(ctx context.Context, c *models.Client) (*models.Client, error)
}

type CarService interface {
	GetBrandByID(ctx context.Context, id string) (*models.Brand, error)
	CreateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error)
	UpdateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error)

	GetModelByID(ctx context.Context, id string) (*models.Model, *models.Brand, error)
	GetModels(ctx context.Context, filters *models.ModelFilters) ([]*models.ModelShort, error)
	CreateModel(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error)
	UpdateModel(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error)

	GetCarByID(ctx context.Context, id string) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
	GetCars(ctx context.Context, filters *models.CarFilters) ([]*models.CarShort, error)
	CreateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
	UpdateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
}

type SupplierService interface {
	GetByID(ctx context.Context, id string) (*models.Supplier, error)
	Create(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
	Update(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
}

type EventService interface {
	AddEvent(ctx context.Context, e *models.Event) error
	Shutdown(ctx context.Context) error
	GetEventByID(ctx context.Context, id string) (*models.Event, error)
	GetEvents(ctx context.Context, filters *models.EventFilters) ([]*models.Event, error)
}
