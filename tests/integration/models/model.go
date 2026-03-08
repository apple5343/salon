package models

import (
	httpModels "salon/internal/transport/http/models"

	"github.com/brianvoe/gofakeit"
)

func GenerateModel(brandID string) *httpModels.Model {
	bodies := make([]string, 0, len(httpModels.BodyTypeMap))
	for _, b := range httpModels.BodyTypeMap {
		bodies = append(bodies, string(b))
	}
	transmissions := make([]string, 0, len(httpModels.TransmissionTypeMap))
	for _, t := range httpModels.TransmissionTypeMap {
		transmissions = append(transmissions, string(t))
	}

	fuels := make([]string, 0, len(httpModels.FuelTypeMap))
	for _, f := range httpModels.FuelTypeMap {
		fuels = append(fuels, string(f))
	}

	drives := make([]string, 0, len(httpModels.DriveTypeMap))
	for _, d := range httpModels.DriveTypeMap {
		drives = append(drives, string(d))
	}

	return &httpModels.Model{
		BrandID:            brandID,
		Name:               gofakeit.Word(),
		Generation:         gofakeit.Word(),
		BodyType:           bodies[gofakeit.Number(0, len(bodies)-1)],
		TransmissionType:   transmissions[gofakeit.Number(0, len(transmissions)-1)],
		FuelType:           fuels[gofakeit.Number(0, len(fuels)-1)],
		EngineDisplacement: gofakeit.Number(5, 1000),
		PowerHP:            gofakeit.Number(5, 1000),
		DriveType:          drives[gofakeit.Number(0, len(drives)-1)],
		BasePrice:          gofakeit.Number(100, 3000),
	}
}

func ModelInternalToModel(m *httpModels.ModelInternalResponse) *httpModels.Model {
	return &httpModels.Model{
		ID:                     m.ID,
		BrandID:                m.Brand.ID,
		Name:                   m.Name,
		Generation:             m.Generation,
		BodyType:               m.BodyType,
		TransmissionType:       m.TransmissionType,
		FuelType:               m.FuelType,
		EngineDisplacement:     m.EngineDisplacement,
		PowerHP:                m.PowerHP,
		DriveType:              m.DriveType,
		BasePrice:              m.BasePrice,
		TechnicCharacteristics: m.TechnicCharacteristics,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}
