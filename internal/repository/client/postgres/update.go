package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	sqlutil "salon/internal/utils/sql"
)

func (r *clientRepository) Update(ctx context.Context, c *service.Client) (*service.Client, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE clients SET 
		full_name = $1, phone = $2, email = $3, password_hash = $4, passport_series = $5, 
		passport_number = $6, passport_issued_by = $7, birth_date = $8, updated_at = $9 WHERE id = $10`,
		c.FullName, c.Phone, c.Email, c.PasswordHash, c.Passport.Series, c.Passport.Number, c.Passport.IssuedBy, c.BirthDate, c.UpdatedAt, c.ID)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.ErrNotFound
	}
	return r.GetByID(ctx, c.ID)
}
