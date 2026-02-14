package models

import (
	"salon/internal/models"
	"time"
)

var BrandOrderByMap = map[string]models.BrandOrderBy{
	"created_at": models.BrandOrderByCreatedAt,
	"updated_at": models.BrandOrderByUpdatedAt,
}

type Brand struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"country_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BrandPublicResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Description string `json:"description"`
}

type BrandInternalResponse struct {
	BrandPublicResponse
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BrandPublicToHttp(b *models.Brand) *BrandPublicResponse {
	return &BrandPublicResponse{
		ID:          b.ID,
		Name:        b.Name,
		CountryCode: b.CountryCode,
		Description: b.Description,
	}
}

func BrandInternalToHttp(b *models.Brand) *BrandInternalResponse {
	return &BrandInternalResponse{
		BrandPublicResponse: *BrandPublicToHttp(b),
		CreatedAt:           b.CreatedAt,
		UpdatedAt:           b.UpdatedAt,
	}
}

func BrandToService(b *Brand) *models.Brand {
	return &models.Brand{
		ID:          b.ID,
		Name:        b.Name,
		CountryCode: b.CountryCode,
		Description: b.Description,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
