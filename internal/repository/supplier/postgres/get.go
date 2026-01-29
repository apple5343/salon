package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *supplierRepository) GetByID(ctx context.Context, id string) (*service.Supplier, error) {
	var s models.Supplier
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM suppliers WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.SupplierToService(&s), nil
}
