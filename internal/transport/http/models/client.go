package models

import (
	"salon/internal/models"
	"time"
)

const TimeLayout = "02.01.2006"

type Client struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Passport  Passport  `json:"passport"`
	BirthDate string    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ClientToHttp(c *models.Client) *Client {
	return &Client{
		ID:        c.ID,
		FullName:  c.FullName,
		Phone:     c.Phone,
		Email:     c.Email,
		Passport:  PassportToHttp(c.Passport),
		BirthDate: c.BirthDate.Format(TimeLayout),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func ClientToService(c *Client) (*models.Client, error) {
	birthDate, err := toTime(c.BirthDate)
	if err != nil {
		return nil, err
	}
	return &models.Client{
		ID:           c.ID,
		FullName:     c.FullName,
		Phone:        c.Phone,
		Email:        c.Email,
		PasswordHash: c.Password,
		Passport:     PassportToService(c.Passport),
		BirthDate:    birthDate,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}, nil
}

func toTime(t string) (time.Time, error) {
	return time.Parse(TimeLayout, t)
}
