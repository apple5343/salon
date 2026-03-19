package postgres

import (
	"context"
	"salon/internal/repository/errors"
	sqlutil "salon/internal/utils/sql"
)

func (r *supplierRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM suppliers WHERE id = $1`, id)
	if err != nil {
		if sqlutil.IsForeignKeyViolationSQL(err) {
			return errors.ErrForeignKey
		}
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}
