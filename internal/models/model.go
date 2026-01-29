package models

import (
	"errors"
	"time"
)

type BodyType string

const (
	BodyTypeSedan       = "sedan"
	BodyTypeHatchback   = "hatchback"
	BodyTypeWagon       = "wagon"
	BodyTypeCrossover   = "crossover"
	BodyTypeSUV         = "suv"
	BodyTypeCoupe       = "coupe"
	BodyTypeConvertible = "convertible"
	BodyTypePickup      = "pickup"
	BodyTypeMinivan     = "minivan"
	BodyTypeVan         = "van"
)

type TransmissionType string

const (
	TransmissionTypeManual      = "manual"
	TransmissionTypeAutomatic   = "automatic"
	TransmissionTypeCVT         = "cvt"
	TransmissionTypeDCT         = "dct"
	TransmissionTypeSingleSpeed = "single-speed"
)

type FuelType string

const (
	FuelTypeGasoline     = "gasoline"
	FuelTypeDiesel       = "diesel"
	FuelTypeElectric     = "electric"
	FuelTypeHybrid       = "hybrid"
	FuelTypePlugInHybrid = "plug-in-hybrid"
	FuelTypeCNG          = "cng"
	FuelTypeHydrogen     = "hydrogen"
)

type DriveType string

const (
	DriveTypeFwd = "fwd"
	DriveTypeAwd = "awd"
	DriveTypeRwd = "rwd"
	DriveType4wd = "4wd"
)

type Model struct {
	ID                     string
	BrandID                string           `validate:"required"`
	Name                   string           `validate:"required"`
	Generation             string           `validate:"required"`
	BodyType               BodyType         `validate:"required"`
	TransmissionType       TransmissionType `validate:"required"`
	FuelType               FuelType         `validate:"required"`
	EngineDisplacement     int              `validate:"required"`
	PowerHP                int              `validate:"required"`
	DriveType              DriveType        `validate:"required"`
	BasePrice              int              `validate:"required"`
	TechnicCharacteristics map[string]interface{}
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (m *Model) BeforeCreate() error {
	if err := m.Validate(); err != nil {
		return err
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func (m *Model) BeforeUpdate() error {
	if err := m.Validate(); err != nil {
		return err
	}
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Model) Validate() error {
	return Validator().Struct(m)
}

type ModelShort struct {
	ID         string
	BrandName  string
	Name       string
	Generation string
	BodyType   BodyType
	DriveType  DriveType
	PowerHP    int
	BasePrice  int
}

type ModelOrderBy string

const (
	ModelOrderByEngineDisplacement = "engine_displacement"
	ModelOrderByPowerHP            = "power_hp"
	ModelOrderByBasePrice          = "base_price"
	ModelOrderByName               = "name"
)

type ModelFilters struct {
	BrandID               *string `validate:"omitempty,uuid"`
	Name                  *string
	Generation            *string
	BodyType              *BodyType
	TransmissionType      *TransmissionType
	FuelType              *FuelType
	MinEngineDisplacement *int `validate:"omitempty,min=0"`
	MaxEngineDisplacement *int `validate:"omitempty,min=0"`
	MinPowerHP            *int `validate:"omitempty,min=0"`
	MaxPowerHP            *int `validate:"omitempty,min=0"`
	DriveType             *DriveType
	MinBasePrice          *int `validate:"omitempty,min=0"`
	MaxBasePrice          *int `validate:"omitempty,min=0"`
	Limit                 *int `validate:"omitempty,min=1"`
	Offset                *int `validate:"omitempty,min=0"`
	OrderBy               *ModelOrderBy
	OrderDirection        *OrderDirection
}

func (f *ModelFilters) Validate() error {
	if err := Validator().Struct(f); err != nil {
		return err
	}
	if f.MinEngineDisplacement != nil && f.MaxEngineDisplacement != nil && *f.MinEngineDisplacement > *f.MaxEngineDisplacement {
		return errors.New("min engine displacement must be less than max engine displacement")
	}
	if f.MinPowerHP != nil && f.MaxPowerHP != nil && *f.MinPowerHP > *f.MaxPowerHP {
		return errors.New("min power hp must be less than max power hp")
	}
	if f.MinBasePrice != nil && f.MaxBasePrice != nil && *f.MinBasePrice > *f.MaxBasePrice {
		return errors.New("min base price must be less than max base price")
	}
	return nil
}
