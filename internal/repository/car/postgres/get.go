package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *carRepository) GetBrandByID(ctx context.Context, id string) (*service.Brand, error) {
	var b models.Brand
	if err := r.db.GetContext(ctx, &b, "SELECT * FROM brands WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.BrandToService(&b), nil
}

func (r *carRepository) GetModelByID(ctx context.Context, id string) (*service.Model, error) {
	var m models.Model
	if err := r.db.GetContext(ctx, &m, "SELECT * FROM models WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.ModelToService(&m), nil
}

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
