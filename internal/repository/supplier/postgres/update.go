package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
)

func (r *supplierRepository) Update(ctx context.Context, s *service.Supplier) (*service.Supplier, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE suppliers SET name = $1, country_code = $2, updated_at = $3 WHERE id = $4`,
		s.Name, s.CountryCode, s.UpdatedAt, s.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
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
	return r.GetByID(ctx, s.ID)
}
