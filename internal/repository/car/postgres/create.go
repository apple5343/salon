package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *carRepository) CreateCar(ctx context.Context, c *service.Car) (*service.Car, error) {
	id := uuid.New().String()
	repoC := models.CarToDatabase(c)
	repoC.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO cars
		(id, model_id, supplier_id, vin, year, color, interior_color, mileage, price, status, options, created_at, updated_at)
		VALUES (:id, :model_id, :supplier_id, :vin, :year, :color, :interior_color, :mileage, :price, :status, :options, :created_at, :updated_at)`,
		repoC)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		} else if sqlutil.IsForeignKeyViolationSQL(err) {
			return nil, errors.ErrForeignKey
		}
		return nil, err
	}
	return r.GetCarByID(ctx, id)
}
