package models

import (
	service "salon/internal/models"
	"time"
)

type Supplier struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	CountryCode string    `json:"country_code" db:"country_code"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func SupplierToDatabase(s *service.Supplier) *Supplier {
	return &Supplier{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func SupplierToService(s *Supplier) *service.Supplier {
	return &service.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
