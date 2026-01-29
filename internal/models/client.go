package models

import (
	"errors"
	"salon/internal/utils/password"
	"time"
)

var ErrPasswordTooShort = errors.New("password must be at least 8 characters")

type Client struct {
	ID           string
	FullName     string `validate:"required"`
	Phone        string `validate:"required,phone"`
	Email        string `validate:"required,email"`
	PasswordHash string // до хеширования сюда приходит пароль из запроса; при Update может быть пустым (не менять)
	Passport     Passport
	BirthDate    time.Time `validate:"required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Client) BeforeCreate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.PasswordHash == "" || len(c.PasswordHash) < 8 {
		return ErrPasswordTooShort
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	hashed, err := password.HashPassword(c.PasswordHash)
	if err != nil {
		return err
	}
	c.PasswordHash = hashed
	return nil
}

func (c *Client) BeforeUpdate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Client) Validate() error {
	return Validator().Struct(c)
}
