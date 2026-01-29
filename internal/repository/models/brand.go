package models

import (
	"salon/internal/models"
	"time"
)

type Brand struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	CountryCode string    `json:"country_code" db:"country_code"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func BrandToDatabase(b *models.Brand) *Brand {
	return &Brand{
		ID:          b.ID,
		Name:        b.Name,
		CountryCode: b.CountryCode,
		Description: b.Description,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
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
