package repository

import (
	"context"
	"salon/internal/models"
)

type EmployeeRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.Employee, error)
	GetByID(ctx context.Context, id string) (*models.Employee, error)
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	Update(ctx context.Context, e *models.Employee) (*models.Employee, error)
}

type ClientRepository interface {
	GetByID(ctx context.Context, id string) (*models.Client, error)
	Create(ctx context.Context, c *models.Client) (*models.Client, error)
	Update(ctx context.Context, c *models.Client) (*models.Client, error)
}

type CarRepository interface {
	GetBrandByID(ctx context.Context, id string) (*models.Brand, error)
	GetBrandsByFilter(ctx context.Context, filter *models.BrandFilters) ([]*models.Brand, error)
	CreateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error)
	UpdateBrand(ctx context.Context, b *models.Brand) (*models.Brand, error)

	GetModelByID(ctx context.Context, id string) (*models.Model, error)
	GetModelsByFilter(ctx context.Context, filter *models.ModelFilters) ([]*models.ModelShort, error)
	CreateModel(ctx context.Context, m *models.Model) (*models.Model, error)
	UpdateModel(ctx context.Context, m *models.Model) (*models.Model, error)

	GetCarByID(ctx context.Context, id string) (*models.Car, error)
	GetCarsByFilter(ctx context.Context, filter *models.CarFilters) ([]*models.CarShort, error)
	CreateCar(ctx context.Context, c *models.Car) (*models.Car, error)
	UpdateCar(ctx context.Context, c *models.Car) (*models.Car, error)
}

type SaleRepository interface {
	GetByID(ctx context.Context, id string) (*models.Sale, error)
	GetSalesByFilter(ctx context.Context, filter *models.SaleFilters) ([]*models.Sale, error)
	Create(ctx context.Context, s *models.Sale) (*models.Sale, error)
	Complete(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
}

type SupplierRepository interface {
	GetByID(ctx context.Context, id string) (*models.Supplier, error)
	GetSuppliersByFilter(ctx context.Context, filter *models.SupplierFilters) ([]*models.Supplier, error)
	Create(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
	Update(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
}

type EventRepository interface {
	Create(ctx context.Context, e *models.Event) (*models.Event, error)
	GetByID(ctx context.Context, id string) (*models.Event, error)
	GetEventsByFilter(ctx context.Context, filter *models.EventFilters) ([]*models.Event, error)
}
