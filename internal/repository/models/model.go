package models

import (
	"salon/internal/models"
	"time"
)

type Model struct {
	ID                     string    `json:"id" db:"id"`
	BrandID                string    `json:"brand_id" db:"brand_id"`
	Name                   string    `json:"name" db:"name"`
	Generation             string    `json:"generation" db:"generation"`
	BodyType               string    `json:"body_type" db:"body_type"`
	TransmissionType       string    `json:"transmission_type" db:"transmission_type"`
	FuelType               string    `json:"fuel_type" db:"fuel_type"`
	EngineDisplacement     int       `json:"engine_displacement" db:"engine_displacement"`
	PowerHP                int       `json:"power_hp" db:"power_hp"`
	DriveType              string    `json:"drive_type" db:"drive_type"`
	BasePrice              int       `json:"base_price" db:"base_price"`
	TechnicCharacteristics JSONB     `json:"technical_characteristics" db:"technical_characteristics"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

func ModelToDatabase(model *models.Model) *Model {
	return &Model{
		ID:                     model.ID,
		BrandID:                model.BrandID,
		Name:                   model.Name,
		Generation:             model.Generation,
		BodyType:               string(model.BodyType),
		TransmissionType:       string(model.TransmissionType),
		FuelType:               string(model.FuelType),
		EngineDisplacement:     model.EngineDisplacement,
		PowerHP:                model.PowerHP,
		DriveType:              string(model.DriveType),
		BasePrice:              model.BasePrice,
		TechnicCharacteristics: JSONB(model.TechnicCharacteristics),
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
	}
}

func ModelToService(model *Model) *models.Model {
	return &models.Model{
		ID:                     model.ID,
		BrandID:                model.BrandID,
		Name:                   model.Name,
		Generation:             model.Generation,
		BodyType:               models.BodyType(model.BodyType),
		TransmissionType:       models.TransmissionType(model.TransmissionType),
		FuelType:               models.FuelType(model.FuelType),
		EngineDisplacement:     model.EngineDisplacement,
		PowerHP:                model.PowerHP,
		DriveType:              models.DriveType(model.DriveType),
		BasePrice:              model.BasePrice,
		TechnicCharacteristics: model.TechnicCharacteristics,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
	}
}

type ModelShort struct {
	ID         string `json:"id" db:"id"`
	BrandName  string `json:"brand_name" db:"brand_name"`
	Name       string `json:"name" db:"name"`
	Generation string `json:"generation" db:"generation"`
	BodyType   string `json:"body_type" db:"body_type"`
	DriveType  string `json:"drive_type" db:"drive_type"`
	PowerHP    int    `json:"power_hp" db:"power_hp"`
	BasePrice  int    `json:"base_price" db:"base_price"`
}

func ModelShortToService(model *ModelShort) *models.ModelShort {
	return &models.ModelShort{
		ID:         model.ID,
		BrandName:  model.BrandName,
		Name:       model.Name,
		Generation: model.Generation,
		BodyType:   models.BodyType(model.BodyType),
		DriveType:  models.DriveType(model.DriveType),
		PowerHP:    model.PowerHP,
		BasePrice:  model.BasePrice,
	}
}
