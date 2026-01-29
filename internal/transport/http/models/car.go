package models

import (
	"salon/internal/models"
	"time"

	"github.com/apple5343/errorx"
	"github.com/shopspring/decimal"
)

var CarStatusType = map[string]models.CarStatus{
	"available": models.CarStatusAvailable,
	"sold":      models.CarStatusSold,
	"pending":   models.CarStatusPending,
	"archived":  models.CarStatusArchived,
	"incoming":  models.CarStatusIncoming,
}

type Car struct {
	ID            string                 `json:"id"`
	ModelID       string                 `json:"model_id"`
	SupplierID    string                 `json:"supplier_id"`
	Vin           string                 `json:"vin"`
	Year          int                    `json:"year"`
	Color         string                 `json:"color"`
	InteriorColor string                 `json:"interior_color"`
	Mileage       int                    `json:"mileage"`
	Price         string                 `json:"price"`
	Status        string                 `json:"status"`
	Options       map[string]interface{} `json:"options"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type CarShort struct {
	ID           string `json:"id"`
	ModelName    string `json:"model_name"`
	SupplierName string `json:"supplier_name"`
	BrandName    string `json:"brand_name"`
	Vin          string `json:"vin"`
	Status       string `json:"status"`
	Price        string `json:"price"`
	Year         int    `json:"year"`
}

type CarPublicResponse struct {
	ID            string                  `json:"id"`
	Model         *ModelPublicResponse    `json:"model"`
	SupplierID    *SupplierPublicResponse `json:"supplier"`
	Vin           string                  `json:"vin"`
	Year          int                     `json:"year"`
	Color         string                  `json:"color"`
	InteriorColor string                  `json:"interior_color"`
	Mileage       int                     `json:"mileage"`
	Price         string                  `json:"price"`
	Status        string                  `json:"status"`
	Options       map[string]interface{}  `json:"options"`
}

type CarInternalResponse struct {
	CarPublicResponse
	Model     *ModelInternalResponse    `json:"model"`
	Supplier  *SupplierInternalResponse `json:"supplier"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func CarShortToHttp(c *models.CarShort) *CarShort {
	return &CarShort{
		ID:           c.ID,
		ModelName:    c.ModelName,
		SupplierName: c.SupplierName,
		BrandName:    c.BrandName,
		Vin:          c.Vin,
		Status:       string(c.Status),
		Price:        c.Price.String(),
		Year:         c.Year,
	}
}

func CarPublicToHttp(c *models.Car, m *models.Model, b *models.Brand, s *models.Supplier) *CarPublicResponse {
	return &CarPublicResponse{
		ID:            c.ID,
		Model:         ModelPublicToHttp(m, b),
		SupplierID:    SupplierPublicToHttp(s),
		Vin:           c.Vin,
		Year:          c.Year,
		Color:         c.Color,
		InteriorColor: c.InteriorColor,
		Mileage:       c.Mileage,
		Price:         c.Price.String(),
		Status:        string(c.Status),
		Options:       c.Options,
	}
}

func CarInternalToHttp(c *models.Car, m *models.Model, b *models.Brand, s *models.Supplier) *CarInternalResponse {
	return &CarInternalResponse{
		CarPublicResponse: *CarPublicToHttp(c, m, b, s),
		Model:             ModelInternalToHttp(m, b),
		Supplier:          SupplierInternalToHttp(s),
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func CarToService(c *Car) (*models.Car, error) {
	status, ok := CarStatusType[c.Status]
	if !ok {
		return nil, errorx.NewError("invalid status", errorx.BadRequest)
	}
	price, err := decimal.NewFromString(c.Price)
	if err != nil {
		return nil, errorx.NewError("invalid price", errorx.BadRequest)
	}
	return &models.Car{
		ID:            c.ID,
		ModelID:       c.ModelID,
		SupplierID:    c.SupplierID,
		Vin:           c.Vin,
		Year:          c.Year,
		Color:         c.Color,
		InteriorColor: c.InteriorColor,
		Mileage:       c.Mileage,
		Price:         price,
		Status:        status,
		Options:       c.Options,
	}, nil
}
