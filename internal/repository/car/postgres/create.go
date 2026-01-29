package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *carRepository) CreateBrand(ctx context.Context, b *service.Brand) (*service.Brand, error) {
	id := uuid.New().String()
	repoB := models.BrandToDatabase(b)
	repoB.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO brands (id, name, country_code, description, created_at, updated_at)
	VALUES (:id, :name, :country_code, :description, :created_at, :updated_at)`, repoB)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	return r.GetBrandByID(ctx, id)
}

func (r *carRepository) CreateModel(ctx context.Context, m *service.Model) (*service.Model, error) {
	id := uuid.New().String()
	repoM := models.ModelToDatabase(m)
	repoM.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO models 
		(id, brand_id, name, generation, body_type, transmission_type, 
		fuel_type, engine_displacement, power_hp, drive_type, base_price, 
		technical_characteristics, created_at, updated_at)
		VALUES (:id, :brand_id, :name, :generation, :body_type, :transmission_type, 
		:fuel_type, :engine_displacement, :power_hp, :drive_type, :base_price, 
		:technical_characteristics, :created_at, :updated_at)`, repoM)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		} else if sqlutil.IsForeignKeyViolationSQL(err) {
			return nil, errors.ErrForeignKey
		}
		return nil, err
	}
	return r.GetModelByID(ctx, id)
}

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
