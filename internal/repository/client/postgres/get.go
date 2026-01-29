package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *clientRepository) GetByID(ctx context.Context, id string) (*service.Client, error) {
	var c models.Client
	if err := r.db.GetContext(ctx, &c, "SELECT * FROM clients WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.ClientToService(&c), nil
}
