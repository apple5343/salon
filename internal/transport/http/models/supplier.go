package models

import (
	"salon/internal/models"
	"time"
)

var SupplierOrderByMap = map[string]models.SupplierOrderBy{
	"created_at": models.SupplierOrderByCreatedAt,
	"updated_at": models.SupplierOrderByUpdatedAt,
}

type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupplierPublicResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
}

type SupplierInternalResponse struct {
	SupplierPublicResponse
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SupplierPublicToHttp(s *models.Supplier) *SupplierPublicResponse {
	return &SupplierPublicResponse{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
	}
}

func SupplierInternalToHttp(s *models.Supplier) *SupplierInternalResponse {
	return &SupplierInternalResponse{
		SupplierPublicResponse: *SupplierPublicToHttp(s),
		CreatedAt:              s.CreatedAt,
		UpdatedAt:              s.UpdatedAt,
	}
}

func SupplierToService(s *Supplier) *models.Supplier {
	return &models.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
