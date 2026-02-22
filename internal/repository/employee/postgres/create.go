package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *employeeRepository) Create(ctx context.Context, e *service.Employee) (*service.Employee, error) {
	repoE := models.EmployeeToDatabase(e)
	id := uuid.New().String()
	repoE.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO employees
		(id, full_name, phone, email, password_hash, passport_series, passport_number, passport_issued_by, role, status, hire_date, created_at, updated_at)
		VALUES (:id, :full_name, :phone, :email, :password_hash, :passport_series, :passport_number, :passport_issued_by, :role, :status, :hire_date, :created_at, :updated_at)`,
		repoE)

	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	return r.GetByID(ctx, id)
}
