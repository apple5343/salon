package service

import (
	"context"
	"salon/internal/models"
	"time"
)

type EmployeeService interface {
	Register(ctx context.Context, e *models.Employee) (*models.Employee, error)
	Hire(ctx context.Context, id string) (*models.Employee, error)
	RegisterAdmin(ctx context.Context, e *models.Employee) (*models.Employee, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	GetRefreshToken(ctx context.Context) (string, error)
	GetAccessToken(ctx context.Context) (string, error)

	GetByID(ctx context.Context, id string) (*models.Employee, error)
	GetEmployees(ctx context.Context, filters *models.EmployeeFilters) ([]*models.Employee, error)
	//TODO stat
	Profile(ctx context.Context) (*models.Employee, error)
	Update(ctx context.Context, e *models.Employee) (*models.Employee, error)
}

type ClientService interface {
	GetByID(ctx context.Context, id string) (*models.Client, error)
	Register(ctx context.Context, c *models.Client) (*models.Client, error)
	//TODO list
	Update(ctx context.Context, c *models.Client) (*models.Client, error)
}

type CarService interface {
	GetCarByID(ctx context.Context, id string) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
	GetCars(ctx context.Context, filters *models.CarFilters) ([]*models.CarShort, error)
	CreateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
	UpdateCar(ctx context.Context, c *models.Car) (*models.Car, *models.Model, *models.Brand, *models.Supplier, error)
}

type BrandService interface {
	GetByID(ctx context.Context, id string) (*models.Brand, error)
	GetBrands(ctx context.Context, filters *models.BrandFilters) ([]*models.Brand, error)
	Create(ctx context.Context, b *models.Brand) (*models.Brand, error)
	Update(ctx context.Context, b *models.Brand) (*models.Brand, error)
	Delete(ctx context.Context, id string) error
}

type ModelService interface {
	GetByID(ctx context.Context, id string) (*models.Model, *models.Brand, error)
	GetModels(ctx context.Context, filters *models.ModelFilters) ([]*models.ModelShort, error)
	Create(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error)
	Update(ctx context.Context, m *models.Model) (*models.Model, *models.Brand, error)
}

type SaleService interface {
	GetByID(ctx context.Context, id string) (*models.Sale, error)
	GetSales(ctx context.Context, filters *models.SaleFilters) ([]*models.Sale, error)
	Create(ctx context.Context, s *models.Sale) (*models.Sale, error)
	Complete(ctx context.Context, id string) (*models.Sale, error)
	Cancel(ctx context.Context, id string) (*models.Sale, error)
}

type SupplierService interface {
	GetByID(ctx context.Context, id string) (*models.Supplier, error)
	GetSuppliers(ctx context.Context, filters *models.SupplierFilters) ([]*models.Supplier, error)
	Create(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
	Update(ctx context.Context, s *models.Supplier) (*models.Supplier, error)
	Delete(ctx context.Context, id string) error
	//TODO stat
}

type EventService interface {
	AddEvent(ctx context.Context, e *models.Event) error
	Shutdown(ctx context.Context) error
	GetEventByID(ctx context.Context, id string) (*models.Event, error)
	GetEvents(ctx context.Context, filters *models.EventFilters) ([]*models.Event, error)
}

type AnalyticsService interface {
	Sales(ctx context.Context, dateFrom, dateTo *time.Time) (*models.SalesAnalytics, error)
	Warehouse(ctx context.Context, turnoverDateFrom, turnoverDateTo *time.Time) (*models.WarehouseAnalytics, error)
	Employee(ctx context.Context, employeeID string, dateFrom, dateTo *time.Time) (*models.EmployeeAnalytics, error)
	Supplier(ctx context.Context, supplierID string, dateFrom, dateTo *time.Time) (*models.SupplierAnalytics, error)
}
