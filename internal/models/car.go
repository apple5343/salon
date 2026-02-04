package models

import (
	"errors"
	"time"

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
	Year          int             `validate:"required,min=1900"`
	Color         string          `validate:"required"`
	InteriorColor string          `validate:"required"`
	Mileage       int             `validate:"required,min=0"`
	Price         decimal.Decimal `validate:"required"`
	Status        CarStatus       `validate:"required"`
	Options       map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c *Car) BeforeCreate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	c.Price = c.Price.Round(2)
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

func (c *Car) BeforeUpdate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	c.Price = c.Price.Round(2)
	c.UpdatedAt = time.Now()
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
	Limit      *int `validate:"omitempty,min=1"`
	Offset     *int `validate:"omitempty,min=0"`
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
