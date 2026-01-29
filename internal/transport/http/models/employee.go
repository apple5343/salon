package models

import (
	"salon/internal/models"
	"time"

	"github.com/apple5343/errorx"
)

var (
	ErrInvalidRole   = errorx.NewError("invalid role", errorx.BadRequest)
	ErrInvalidStatus = errorx.NewError("invalid status", errorx.BadRequest)
)

type Employee struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Passport  Passport  `json:"passport"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	HireDate  time.Time `json:"hire_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func EmployeeToHttp(e *models.Employee) *Employee {
	return &Employee{
		ID:        e.ID,
		FullName:  e.FullName,
		Phone:     e.Phone,
		Email:     e.Email,
		Passport:  PassportToHttp(e.Passport),
		Role:      string(e.Role),
		Status:    string(e.Status),
		HireDate:  e.HireDate,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func EmployeeToService(e *Employee) (*models.Employee, error) {
	role, err := ToServiceRole(e.Role)
	if err != nil {
		return nil, err
	}

	status, err := ToServiceStatus(e.Status)
	if err != nil {
		return nil, err
	}
	return &models.Employee{
		ID:           e.ID,
		FullName:     e.FullName,
		Phone:        e.Phone,
		Email:        e.Email,
		PasswordHash: e.Password,
		Passport:     PassportToService(e.Passport),
		Role:         role,
		Status:       status,
		HireDate:     e.HireDate,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}, nil
}

func ToServiceRole(role string) (models.EmployeeRole, error) {
	switch role {
	case "admin":
		return models.EmployeeRoleAdmin, nil
	case "manager":
		return models.EmployeeRoleManager, nil
	default:
		return "", ErrInvalidRole
	}
}

func ToServiceStatus(status string) (models.EmployeeStatus, error) {
	switch status {
	case "active":
		return models.EmployeeStatusActive, nil
	case "inactive":
		return models.EmployeeStatusInactive, nil
	default:
		return "", ErrInvalidStatus
	}
}
