package models

import (
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
)

func GenerateBrand() *httpModels.Brand {
	return &httpModels.Brand{
		Name:        gofakeit.Name(),
		CountryCode: CountyCodes[gofakeit.Number(0, len(CountyCodes)-1)],
		Description: gofakeit.Word(),
	}
}

func CopyBrand(brand *httpModels.Brand) *httpModels.Brand {
	return &httpModels.Brand{
		ID:          brand.ID,
		Name:        brand.Name,
		CountryCode: brand.CountryCode,
		Description: brand.Description,
		CreatedAt:   brand.CreatedAt,
		UpdatedAt:   brand.UpdatedAt,
	}
}
