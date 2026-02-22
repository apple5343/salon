package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	sqlutil "salon/internal/utils/sql"
)

func (r *brandRepository) Update(ctx context.Context, b *service.Brand) (*service.Brand, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE brands SET 
		name = $1, country_code = $2, description = $3, updated_at = $4 WHERE id = $5`,
		b.Name, b.CountryCode, b.Description, b.UpdatedAt, b.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		} else if sqlutil.IsUniqueViolationSQL(err) {
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
	return r.GetByID(ctx, b.ID)
}
