package models

import (
	"salon/internal/utils/password"
	"time"
)

type EmployeeRole string
type EmployeeStatus string

const (
	EmployeeRoleAdmin   EmployeeRole = "admin"
	EmployeeRoleManager EmployeeRole = "manager"

	EmployeeStatusActive   EmployeeStatus = "active"
	EmployeeStatusInactive EmployeeStatus = "inactive"
)

type Employee struct {
	ID           string
	FullName     string `validate:"required"`
	Phone        string `validate:"required,phone"`
	Email        string `validate:"required,email"`
	PasswordHash string
	Passport     Passport
	Role         EmployeeRole   `validate:"required"`
	Status       EmployeeStatus `validate:"required"`
	HireDate     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (e *Employee) BeforeCreate() error {
	if err := e.Validate(); err != nil {
		return err
	}
	if e.PasswordHash == "" || len(e.PasswordHash) < 8 {
		return ErrPasswordTooShort
	}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	passwordHash, err := password.HashPassword(e.PasswordHash)
	if err != nil {
		return err
	}
	e.PasswordHash = passwordHash
	return nil
}

func (e *Employee) BeforeUpdate() error {
	if err := e.Validate(); err != nil {
		return err
	}
	if e.PasswordHash != "" && len(e.PasswordHash) < 8 {
		return ErrPasswordTooShort
	}
	now := time.Now()
	e.UpdatedAt = now
	return nil
}

func (e *Employee) Validate() error {
	return Validator().Struct(e)
}
