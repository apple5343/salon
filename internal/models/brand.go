package models

import "time"

type Brand struct {
	ID          string
	Name        string `validate:"required"`
	CountryCode string `validate:"required,iso3166_1_alpha2"`
	Description string `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *Brand) BeforeCreate() error {
	if err := b.Validate(); err != nil {
		return err
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *Brand) BeforeUpdate() error {
	if err := b.Validate(); err != nil {
		return err
	}
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Brand) Validate() error {
	return Validator().Struct(b)
}

type BrandOrderBy string

const (
	BrandOrderByCreatedAt BrandOrderBy = "created_at"
	BrandOrderByUpdatedAt BrandOrderBy = "updated_at"
)

type BrandFilters struct {
	Name        *string `validate:"omitempty,min=1"`
	CountryCode *string `validate:"omitempty,iso3166_1_alpha2"`
	BaseList
	OrderBy *BrandOrderBy
}

func (f *BrandFilters) Validate() error {
	return Validator().Struct(f)
}
