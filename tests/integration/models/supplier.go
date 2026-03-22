package models

import (
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
)

var CountyCodes = []string{"US", "CA", "GB", "DE", "IT", "RU"}

func GenerateSupplier() *httpModels.Supplier {
	return &httpModels.Supplier{
		Name:        gofakeit.Name(),
		CountryCode: CountyCodes[gofakeit.Number(0, len(CountyCodes)-1)],
	}
}

func CopySupplier(s *httpModels.Supplier) *httpModels.Supplier {
	return &httpModels.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func SupplierInternalToSupplier(s *httpModels.SupplierInternalResponse) *httpModels.Supplier {
	return &httpModels.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		CountryCode: s.CountryCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
