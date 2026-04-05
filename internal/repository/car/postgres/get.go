package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *carRepository) GetCarByID(ctx context.Context, id string) (*service.Car, error) {
	var c models.Car
	if err := r.db.GetContext(ctx, &c, "SELECT * FROM cars WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.CarToService(&c), nil
}
