package models

import (
	"salon/internal/models"
	"time"
)

type Client struct {
	ID               string    `json:"id" db:"id"`
	FullName         string    `json:"full_name" db:"full_name"`
	Phone            string    `json:"phone" db:"phone"`
	Email            string    `json:"email" db:"email"`
	PasswordHash     string    `json:"password_hash" db:"password_hash"`
	PassportSeries   string    `json:"passport_series" db:"passport_series"`
	PassportNumber   string    `json:"passport_number" db:"passport_number"`
	PassportIssuedBy string    `json:"passport_issued_by" db:"passport_issued_by"`
	BirthDate        time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func ClientToDatabase(c *models.Client) *Client {
	return &Client{
		ID:               c.ID,
		FullName:         c.FullName,
		Phone:            c.Phone,
		Email:            c.Email,
		PasswordHash:     c.PasswordHash,
		PassportSeries:   c.Passport.Series,
		PassportNumber:   c.Passport.Number,
		PassportIssuedBy: c.Passport.IssuedBy,
		BirthDate:        c.BirthDate,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func ClientToService(c *Client) *models.Client {
	return &models.Client{
		ID:           c.ID,
		FullName:     c.FullName,
		Phone:        c.Phone,
		Email:        c.Email,
		PasswordHash: c.PasswordHash,
		Passport: models.Passport{
			Series:   c.PassportSeries,
			Number:   c.PassportNumber,
			IssuedBy: c.PassportIssuedBy,
		},
		BirthDate: c.BirthDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
