package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	sqlutil "salon/internal/utils/sql"
)

func (r *employeeRepository) Update(ctx context.Context, e *service.Employee) (*service.Employee, error) {
	res, err := r.db.DB.ExecContext(ctx, `UPDATE employees SET full_name = $1, phone = $2, email = $3, password_hash = $4,
		passport_series = $5, passport_number = $6, passport_issued_by = $7, updated_at = $8 WHERE id = $9`,
		e.FullName, e.Phone, e.Email, e.PasswordHash, e.Passport.Series, e.Passport.Number, e.Passport.IssuedBy, e.UpdatedAt, e.ID)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.ErrNotFound
	}
	return r.GetByID(ctx, e.ID)
}
