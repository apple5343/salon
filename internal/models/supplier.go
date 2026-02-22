package models

import (
	"salon/pkg/clock"
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	ID          string
	Name        string `validate:"required"`
	CountryCode string `validate:"required,iso3166_1_alpha2"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Supplier) BeforeCreate(c clock.Clock) error {
	if err := s.Validate(); err != nil {
		return err
	}
	now := c.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

func (s *Supplier) BeforeUpdate(c clock.Clock) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if _, err := uuid.Parse(s.ID); err != nil {
		return ErrInvalidID
	}
	s.UpdatedAt = c.Now()
	return nil
}

func (s *Supplier) Validate() error {
	return Validator().Struct(s)
}

type SupplierOrderBy string

const (
	SupplierOrderByCreatedAt SupplierOrderBy = "created_at"
	SupplierOrderByUpdatedAt SupplierOrderBy = "updated_at"
)

type SupplierFilters struct {
	Name        *string `validate:"omitempty,min=1"`
	CountryCode *string `validate:"omitempty,iso3166_1_alpha2"`
	BaseList
	OrderBy *SupplierOrderBy
}

func (f *SupplierFilters) Validate() error {
	return Validator().Struct(f)
}
