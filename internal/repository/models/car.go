package models

import (
	"salon/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

type Car struct {
	ID            string          `json:"id" db:"id"`
	ModelID       string          `json:"model_id" db:"model_id"`
	SupplierID    string          `json:"supplier_id" db:"supplier_id"`
	Vin           string          `json:"vin" db:"vin"`
	Year          int             `json:"year" db:"year"`
	Color         string          `json:"color" db:"color"`
	InteriorColor string          `json:"interior_color" db:"interior_color"`
	Mileage       int             `json:"mileage" db:"mileage"`
	Price         decimal.Decimal `json:"price" db:"price"`
	Status        string          `json:"status" db:"status"`
	Options       JSONB           `json:"options" db:"options"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

func CarToDatabase(c *models.Car) *Car {
	return &Car{
		ID:            c.ID,
		ModelID:       c.ModelID,
		SupplierID:    c.SupplierID,
		Vin:           c.Vin,
		Year:          c.Year,
		Color:         c.Color,
		InteriorColor: c.InteriorColor,
		Mileage:       c.Mileage,
		Price:         c.Price,
		Status:        string(c.Status),
		Options:       JSONB(c.Options),
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func CarToService(c *Car) *models.Car {
	return &models.Car{
		ID:            c.ID,
		ModelID:       c.ModelID,
		SupplierID:    c.SupplierID,
		Vin:           c.Vin,
		Year:          c.Year,
		Color:         c.Color,
		InteriorColor: c.InteriorColor,
		Mileage:       c.Mileage,
		Price:         c.Price,
		Status:        models.CarStatus(c.Status),
		Options:       c.Options,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

type CarShort struct {
	ID           string          `json:"id" db:"id"`
	ModelName    string          `json:"model_name" db:"model_name"`
	SupplierName string          `json:"supplier_name" db:"supplier_name"`
	BrandName    string          `json:"brand_name" db:"brand_name"`
	Vin          string          `json:"vin" db:"vin"`
	Status       string          `json:"status" db:"status"`
	Price        decimal.Decimal `json:"price" db:"price"`
	Year         int             `json:"year" db:"year"`
}

func CarShortToService(c *CarShort) *models.CarShort {
	return &models.CarShort{
		ID:           c.ID,
		ModelName:    c.ModelName,
		SupplierName: c.SupplierName,
		BrandName:    c.BrandName,
		Vin:          c.Vin,
		Status:       models.CarStatus(c.Status),
		Price:        c.Price,
		Year:         c.Year,
	}
}
