package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*service.Employee, error) {
	var e models.Employee
	if err := r.db.GetContext(ctx, &e, "SELECT * FROM employees WHERE email = $1", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.EmployeeToService(&e), nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id string) (*service.Employee, error) {
	var e models.Employee
	if err := r.db.GetContext(ctx, &e, "SELECT * FROM employees WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.EmployeeToService(&e), nil
}
