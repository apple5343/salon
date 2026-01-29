package models

import "time"

type Supplier struct {
	ID          string
	Name        string `validate:"required"`
	CountryCode string `validate:"required,iso3166_1_alpha2"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Supplier) BeforeCreate() error {
	if err := s.Validate(); err != nil {
		return err
	}
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

func (s *Supplier) BeforeUpdate() error {
	if err := s.Validate(); err != nil {
		return err
	}
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Supplier) Validate() error {
	return Validator().Struct(s)
}
