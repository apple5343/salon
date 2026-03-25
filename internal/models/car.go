package models

import (
	"errors"
	"salon/pkg/clock"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CarStatus string

var (
	CarStatusAvailable CarStatus = "available"
	CarStatusSold      CarStatus = "sold"
	CarStatusPending   CarStatus = "pending"
	CarStatusArchived  CarStatus = "archived"
	CarStatusIncoming  CarStatus = "incoming"
	CarStatusBooked    CarStatus = "booked"
)

type Car struct {
	ID            string
	ModelID       string          `validate:"required,uuid"`
	SupplierID    string          `validate:"required,uuid"`
	Vin           string          `validate:"required,len=17"`
	Year          int             `validate:"required,min=1900,max=2050"`
	Color         string          `validate:"required"`
	InteriorColor string          `validate:"required"`
	Mileage       int             `validate:"min=0"`
	Price         decimal.Decimal `validate:"required"`
	Status        CarStatus       `validate:"required"`
	Options       map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c *Car) BeforeCreate(cl clock.Clock) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.Price.IsNegative() {
		return errors.New("price cannot be negative")
	}
	c.Price = c.Price.Round(2)
	now := cl.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

func (c *Car) BeforeUpdate(cl clock.Clock) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if _, err := uuid.Parse(c.ID); err != nil {
		return ErrInvalidID
	}
	if c.Price.IsNegative() {
		return errors.New("price cannot be negative")
	}
	c.Price = c.Price.Round(2)
	c.UpdatedAt = cl.Now()
	return nil
}

func (c *Car) Validate() error {
	return Validator().Struct(c)
}

type CarShort struct {
	ID           string
	ModelName    string
	SupplierName string
	BrandName    string
	Vin          string
	Status       CarStatus
	Price        decimal.Decimal
	Year         int
}

type CarOrderBy string

const (
	CarOrderByYear      CarOrderBy = "year"
	CarOrderByPrice     CarOrderBy = "price"
	CarOrderByMile      CarOrderBy = "mileage"
	CarOrderByCreatedAt CarOrderBy = "created_at"
	CarOrderByUpdatedAt CarOrderBy = "updated_at"
)

type CarFilters struct {
	SupplierID *string `validate:"omitempty,uuid"`
	ModelID    *string `validate:"omitempty,uuid"`
	BrandID    *string `validate:"omitempty,uuid"`
	MinYear    *int    `validate:"omitempty,min=1900"`
	MaxYear    *int    `validate:"omitempty,min=1900"`
	Color      *string
	Status     *CarStatus
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	MinMileage *int `validate:"omitempty,min=0"`
	MaxMileage *int `validate:"omitempty,min=0"`
	OrderBy    *CarOrderBy
	BaseList
}

func (f *CarFilters) Validate() error {
	if err := Validator().Struct(f); err != nil {
		return err
	}
	if f.MinPrice != nil {
		*f.MinPrice = f.MinPrice.Round(2)
	}
	if f.MaxPrice != nil {
		*f.MaxPrice = f.MaxPrice.Round(2)
	}
	if f.MinMileage != nil && f.MaxMileage != nil && *f.MinMileage > *f.MaxMileage {
		return errors.New("min mileage must be less than max mileage")
	}
	if f.MinYear != nil && f.MaxYear != nil && *f.MinYear > *f.MaxYear {
		return errors.New("min year must be less than max year")
	}
	return nil
}
