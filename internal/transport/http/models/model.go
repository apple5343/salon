package models

import (
	"errors"
	"salon/internal/models"
	"time"
)

var BodyTypeMap = map[string]models.BodyType{
	"sedan":       models.BodyTypeSedan,
	"coupe":       models.BodyTypeCoupe,
	"hatchback":   models.BodyTypeHatchback,
	"minivan":     models.BodyTypeMinivan,
	"pickup":      models.BodyTypePickup,
	"suv":         models.BodyTypeSUV,
	"van":         models.BodyTypeVan,
	"wagon":       models.BodyTypeWagon,
	"convertible": models.BodyTypeConvertible,
	"crossover":   models.BodyTypeCrossover,
}

var TransmissionTypeMap = map[string]models.TransmissionType{
	"automatic":         models.TransmissionTypeAutomatic,
	"manual":       models.TransmissionTypeManual,
	"cvt":          models.TransmissionTypeCVT,
	"dct":          models.TransmissionTypeDCT,
	"single-speed": models.TransmissionTypeSingleSpeed,
}

var FuelTypeMap = map[string]models.FuelType{
	"gasoline":       models.FuelTypeGasoline,
	"diesel":         models.FuelTypeDiesel,
	"electric":       models.FuelTypeElectric,
	"hybrid":         models.FuelTypeHybrid,
	"cng":            models.FuelTypeCNG,
	"hydrogen":       models.FuelTypeHydrogen,
	"plug-in-hybrid": models.FuelTypePlugInHybrid,
}

var DriveTypeMap = map[string]models.DriveType{
	"awd": models.DriveTypeAwd,
	"fwd": models.DriveTypeFwd,
	"rwd": models.DriveTypeRwd,
	"4wd": models.DriveType4wd,
}

var ModelOrderByMap = map[string]models.ModelOrderBy{
	"name":                models.ModelOrderByName,
	"base_price":          models.ModelOrderByBasePrice,
	"engine_displacement": models.ModelOrderByEngineDisplacement,
	"power_hp":            models.ModelOrderByPowerHP, //TODO add created_at & updated_at
}

type Model struct {
	ID                     string                 `json:"id"`
	BrandID                string                 `json:"brand_id"`
	Name                   string                 `json:"name"`
	Generation             string                 `json:"generation"`
	BodyType               string                 `json:"body_type"`
	TransmissionType       string                 `json:"transmission_type"`
	FuelType               string                 `json:"fuel_type"`
	EngineDisplacement     int                    `json:"engine_displacement"`
	PowerHP                int                    `json:"power_hp"`
	DriveType              string                 `json:"drive_type"`
	BasePrice              int                    `json:"base_price"`
	TechnicCharacteristics map[string]interface{} `json:"technic_characteristics"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type ModelShort struct {
	ID         string `json:"id"`
	BrandName  string `json:"brand_name"`
	Name       string `json:"name"`
	Generation string `json:"generation"`
	BodyType   string `json:"body_type"`
	DriveType  string `json:"drive_type"`
	PowerHP    int    `json:"power_hp"`
	BasePrice  int    `json:"base_price"`
}

type ModelPublicResponse struct {
	ID                     string                 `json:"id"`
	Brand                  *BrandPublicResponse   `json:"brand"`
	Name                   string                 `json:"name"`
	Generation             string                 `json:"generation"`
	BodyType               string                 `json:"body_type"`
	TransmissionType       string                 `json:"transmission_type"`
	FuelType               string                 `json:"fuel_type"`
	EngineDisplacement     int                    `json:"engine_displacement"`
	PowerHP                int                    `json:"power_hp"`
	DriveType              string                 `json:"drive_type"`
	BasePrice              int                    `json:"base_price"`
	TechnicCharacteristics map[string]interface{} `json:"technic_characteristics"`
}

type ModelInternalResponse struct {
	ModelPublicResponse
	Brand     *BrandInternalResponse `json:"brand"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func ModelPublicToHttp(model *models.Model, brand *models.Brand) *ModelPublicResponse {
	return &ModelPublicResponse{
		ID:                     model.ID,
		Brand:                  BrandPublicToHttp(brand),
		Name:                   model.Name,
		Generation:             model.Generation,
		BodyType:               string(model.BodyType),
		TransmissionType:       string(model.TransmissionType),
		FuelType:               string(model.FuelType),
		EngineDisplacement:     model.EngineDisplacement,
		PowerHP:                model.PowerHP,
		DriveType:              string(model.DriveType),
		BasePrice:              model.BasePrice,
		TechnicCharacteristics: model.TechnicCharacteristics,
	}
}

func ModelInternalToHttp(model *models.Model, brand *models.Brand) *ModelInternalResponse {
	return &ModelInternalResponse{
		ModelPublicResponse: *ModelPublicToHttp(model, brand),
		Brand:               BrandInternalToHttp(brand),
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
	}
}

func ModelShortToHttp(model *models.ModelShort) *ModelShort {
	return &ModelShort{
		ID:         model.ID,
		BrandName:  model.BrandName,
		Name:       model.Name,
		Generation: model.Generation,
		BodyType:   string(model.BodyType),
		DriveType:  string(model.DriveType),
		PowerHP:    model.PowerHP,
		BasePrice:  model.BasePrice,
	}
}

func ModelToService(model *Model) (*models.Model, error) {
	bodyType, ok := BodyTypeMap[model.BodyType]
	if !ok {
		return nil, errors.New("invalid body type")
	}

	transmissionType, ok := TransmissionTypeMap[model.TransmissionType]
	if !ok {
		return nil, errors.New("invalid transmission type")
	}

	fuelType, ok := FuelTypeMap[model.FuelType]
	if !ok {
		return nil, errors.New("invalid fuel type")
	}

	driveType, ok := DriveTypeMap[model.DriveType]
	if !ok {
		return nil, errors.New("invalid drive type")
	}

	return &models.Model{
		ID:                     model.ID,
		BrandID:                model.BrandID,
		Name:                   model.Name,
		Generation:             model.Generation,
		BodyType:               bodyType,
		TransmissionType:       transmissionType,
		FuelType:               fuelType,
		EngineDisplacement:     model.EngineDisplacement,
		PowerHP:                model.PowerHP,
		DriveType:              driveType,
		BasePrice:              model.BasePrice,
		TechnicCharacteristics: model.TechnicCharacteristics,
	}, nil
}
