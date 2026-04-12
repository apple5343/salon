package models

import (
	"database/sql"
	"salon/internal/models"
	"time"
)

type Employee struct {
	ID               string       `json:"id" db:"id"`
	FullName         string       `json:"full_name" db:"full_name"`
	Phone            string       `json:"phone" db:"phone"`
	Email            string       `json:"email" db:"email"`
	PasswordHash     string       `json:"password_hash" db:"password_hash"`
	PassportSeries   string       `json:"passport_series" db:"passport_series"`
	PassportNumber   string       `json:"passport_number" db:"passport_number"`
	PassportIssuedBy string       `json:"passport_issued_by" db:"passport_issued_by"`
	Role             string       `json:"role" db:"role"`
	Status           string       `json:"status" db:"status"`
	HireDate         sql.NullTime `json:"hire_date" db:"hire_date"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
}

func EmployeeToDatabase(e *models.Employee) *Employee {
	var hireDate sql.NullTime
	if e.HireDate != nil {
		hireDate = sql.NullTime{Time: *e.HireDate, Valid: true}
	}
	return &Employee{
		ID:               e.ID,
		FullName:         e.FullName,
		Phone:            e.Phone,
		Email:            e.Email,
		PasswordHash:     e.PasswordHash,
		PassportSeries:   e.Passport.Series,
		PassportNumber:   e.Passport.Number,
		PassportIssuedBy: e.Passport.IssuedBy,
		Role:             string(e.Role),
		Status:           string(e.Status),
		HireDate:         hireDate,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func EmployeeToService(e *Employee) *models.Employee {
	var hireDate *time.Time
	if e.HireDate.Valid {
		hireDate = &e.HireDate.Time
	}
	return &models.Employee{
		ID:           e.ID,
		FullName:     e.FullName,
		Phone:        e.Phone,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		Passport: models.Passport{
			Series:   e.PassportSeries,
			Number:   e.PassportNumber,
			IssuedBy: e.PassportIssuedBy,
		},
		Role:      ToServiceRole(e.Role),
		Status:    ToServiceStatus(e.Status),
		HireDate:  hireDate,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToServiceRole(role string) models.EmployeeRole {
	switch role {
	case "admin":
		return models.EmployeeRoleAdmin
	case "manager":
		return models.EmployeeRoleManager
	default:
		return ""
	}
}

func ToServiceStatus(status string) models.EmployeeStatus {
	switch status {
	case "active":
		return models.EmployeeStatusActive
	case "inactive":
		return models.EmployeeStatusInactive
	default:
		return ""
	}
}
