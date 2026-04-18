package models

import (
	"errors"
	"salon/pkg/clock"
	"time"

	"github.com/shopspring/decimal"
)

type PaymentType string

const (
	PaymentTypeCash PaymentType = "cash"
	PaymentTypeCard PaymentType = "card"
)

type SaleStatus string

const (
	SaleStatusPending   SaleStatus = "pending"
	SaleStatusCompleted SaleStatus = "completed"
	SaleStatusCanceled  SaleStatus = "canceled"
)

type Sale struct {
	ID              string
	CarID           string `validate:"required,uuid"`
	ClientID        string `validate:"required,uuid"`
	EmployeeID      string `validate:"required,uuid"`
	SaleDate        time.Time
	OriginPrice     decimal.Decimal
	DiscountAmount  decimal.Decimal
	DiscountPercent decimal.Decimal
	FinalPrice      decimal.Decimal
	PaymentType     PaymentType
	Status          SaleStatus
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *Sale) BeforeCreate(c clock.Clock) error {
	if err := s.Validate(); err != nil {
		return err
	}
	now := c.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	s.SaleDate = now
	return nil
}

func (s *Sale) BeforeUpdate(c clock.Clock) error {
	if err := s.Validate(); err != nil {
		return err
	}
	s.UpdatedAt = c.Now()
	return nil
}

func (s *Sale) Validate() error {
	return Validator().Struct(s)
}

type SaleFilters struct {
	CarID         *string `validate:"omitempty,uuid"`
	ClientID      *string `validate:"omitempty,uuid"`
	EmployeeID    *string `validate:"omitempty,uuid"`
	FinalPriceMin *decimal.Decimal
	FinalPriceMax *decimal.Decimal
	Status        *SaleStatus
	PaymentType   *PaymentType
	DateFrom      *time.Time
	DateTo        *time.Time
	OrderBy       *SaleOrderBy
	BaseList
}

type SaleOrderBy string

const (
	SaleOrderByFinalPrice SaleOrderBy = "final_price"
	SaleOrderByDate       SaleOrderBy = "date"
)

func (s *SaleFilters) Validate() error {
	if err := Validator().Struct(s); err != nil {
		return err
	}
	if s.FinalPriceMin != nil && s.FinalPriceMax != nil && s.FinalPriceMin.GreaterThan(*s.FinalPriceMax) {
		return errors.New("final price min must be less or equal than final price max")
	}
	if s.DateFrom != nil && s.DateTo != nil && s.DateFrom.After(*s.DateTo) {
		return errors.New("date from must be less or equal than date to")
	}
	return nil
}
