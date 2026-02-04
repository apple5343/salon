package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *saleRepository) GetByID(ctx context.Context, id string) (*service.Sale, error) {
	var s models.Sale
	if err := r.db.GetContext(ctx, &s, "SELECT * FROM sales WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.SaleToService(&s), nil
}
