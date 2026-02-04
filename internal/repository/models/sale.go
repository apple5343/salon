package models

import (
	"salon/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

type Sale struct {
	ID              string          `json:"id" db:"id"`
	CarID           string          `json:"car_id" db:"car_id"`
	ClientID        string          `json:"client_id" db:"client_id"`
	EmployeeID      string          `json:"employee_id" db:"employee_id"`
	SaleDate        time.Time       `json:"sale_date" db:"sale_date"`
	OriginPrice     decimal.Decimal `json:"original_price" db:"original_price"`
	DiscountAmount  decimal.Decimal `json:"discount_amount" db:"discount_amount"`
	DiscountPercent decimal.Decimal `json:"discount_percent" db:"discount_percent"`
	FinalPrice      decimal.Decimal `json:"final_price" db:"final_price"`
	PaymentType     string          `json:"payment_type" db:"payment_type"`
	Status          string          `json:"status" db:"status"`
	Notes           string          `json:"notes" db:"notes"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

func SaleToDatabase(sale *models.Sale) *Sale {
	return &Sale{
		ID:              sale.ID,
		CarID:           sale.CarID,
		ClientID:        sale.ClientID,
		EmployeeID:      sale.EmployeeID,
		SaleDate:        sale.SaleDate,
		OriginPrice:     sale.OriginPrice,
		DiscountAmount:  sale.DiscountAmount,
		DiscountPercent: sale.DiscountPercent,
		FinalPrice:      sale.FinalPrice,
		PaymentType:     string(sale.PaymentType),
		Status:          string(sale.Status),
		Notes:           sale.Notes,
		CreatedAt:       sale.CreatedAt,
		UpdatedAt:       sale.UpdatedAt,
	}
}

func SaleToService(sale *Sale) *models.Sale {
	return &models.Sale{
		ID:              sale.ID,
		CarID:           sale.CarID,
		ClientID:        sale.ClientID,
		EmployeeID:      sale.EmployeeID,
		SaleDate:        sale.SaleDate,
		OriginPrice:     sale.OriginPrice,
		DiscountAmount:  sale.DiscountAmount,
		DiscountPercent: sale.DiscountPercent,
		FinalPrice:      sale.FinalPrice,
		PaymentType:     models.PaymentType(sale.PaymentType),
		Status:          models.SaleStatus(sale.Status),
		Notes:           sale.Notes,
		CreatedAt:       sale.CreatedAt,
		UpdatedAt:       sale.UpdatedAt,
	}
}
