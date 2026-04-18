package repository

import (
	"context"
	"salon/internal/models"
	"time"
)

type EmployeeRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.Employee, error)
	GetByID(ctx context.Context, id string) (*models.Employee, error)
	GetEmployeesByFilter(ctx context.Context, filter *models.EmployeeFilters) ([]*models.Employee, error)
	Create(ctx context.Context, e *models.Employee) (*models.Employee, error)
	Update(ctx context.Context, e *models.Employee) (*models.Employee, error)
}

type ClientRepository interface {
	GetByID(ctx context.Context, id string) (*models.Client, error)
	Create(ctx context.Context, c *models.Client) (*models.Client, error)
	Update(ctx context.Context, c *models.Client) (*models.Client, error)
}

type CarRepository interface {
	GetCarByID(ctx context.Context, id string) (*models.Car, error)
	GetCarsByFilter(ctx context.Context, filter *models.CarFilters) ([]*models.CarShort, error)
	CreateCar(ctx context.Context, c *models.Car) (*models.Car, error)
	UpdateCar(ctx context.Context, c *models.Car) (*models.Car, error)
}

type BrandRepository interface {
	GetByID(ctx context.Context, id string) (*models.Brand, error)
	GetBrandsByFilter(ctx context.Context, filter *models.BrandFilters) ([]*models.Brand, error)
	Create(ctx context.Context, b *models.Brand) (*models.Brand, error)
	Update(ctx context.Context, b *models.Brand) (*models.Brand, error)
	Delete(ctx context.Context, id string) error
}

type BrandCache interface {
	GetByID(ctx context.Context, id string) (*models.Brand, error)
	SetByID(ctx context.Context, b *models.Brand, ttl time.Duration) error
	DeleteByID(ctx context.Context, id string) error
}

type ModelRepository interface {
	GetByID(ctx context.Context, id string) (*models.Model, error)
	GetModelsByFilter(ctx context.Context, filter *models.ModelFilters) ([]*models.ModelShort, error)
	Create(ctx context.Context, m *models.Model) (*models.Model, error)
	Update(ctx context.Context, m *models.Model) (*models.Model, error)
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
	Delete(ctx context.Context, id string) error
}

type SupplierCache interface {
	GetByID(ctx context.Context, id string) (*models.Supplier, error)
	SetByID(ctx context.Context, s *models.Supplier, ttl time.Duration) error
	DeleteByID(ctx context.Context, id string) error
}

type EventRepository interface {
	Create(ctx context.Context, e *models.Event) (*models.Event, error)
	GetByID(ctx context.Context, id string) (*models.Event, error)
	GetEventsByFilter(ctx context.Context, filter *models.EventFilters) ([]*models.Event, error)
}

type AnalyticsRepository interface {
	GetSalesAnalytics(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SalesAnalytics, error)
}
