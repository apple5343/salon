package models

import "time"

type Client struct {
	ID           string
	FullName     string `validate:"required"`
	Phone        string `validate:"required,phone"`
	Email        string `validate:"required,email"`
	PasswordHash string
	Passport     Passport
	BirthDate    time.Time `validate:"required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Client) BeforeCreate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
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
