package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"
)

func (r *carRepository) UpdateBrand(ctx context.Context, b *service.Brand) (*service.Brand, error) {
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
	return r.GetBrandByID(ctx, b.ID)
}

func (r *carRepository) UpdateModel(ctx context.Context, m *service.Model) (*service.Model, error) {
	repoM := models.ModelToDatabase(m)
	res, err := r.db.ExecContext(ctx, `UPDATE models SET 
		brand_id = $1, name = $2, generation = $3, body_type = $4, transmission_type = $5,
		fuel_type = $6, engine_displacement = $7, power_hp = $8, drive_type = $9, base_price = $10,
		technical_characteristics = $11, updated_at = $12 WHERE id = $13`,
		repoM.BrandID, repoM.Name, repoM.Generation, repoM.BodyType, repoM.TransmissionType,
		repoM.FuelType, repoM.EngineDisplacement, repoM.PowerHP, repoM.DriveType, repoM.BasePrice,
		repoM.TechnicCharacteristics, repoM.UpdatedAt, repoM.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		} else if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		} else if sqlutil.IsForeignKeyViolationSQL(err) {
			return nil, errors.ErrForeignKey
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
	return r.GetModelByID(ctx, m.ID)
}

func (r *carRepository) UpdateCar(ctx context.Context, c *service.Car) (*service.Car, error) {
	repoC := models.CarToDatabase(c)
	res, err := r.db.ExecContext(ctx, `UPDATE cars SET 
		model_id = $1, supplier_id = $2, vin = $3, year = $4, color = $5, interior_color = $6,
		mileage = $7, price = $8, status = $9, options = $10, updated_at = $11 WHERE id = $12`,
		repoC.ModelID, repoC.SupplierID, repoC.Vin, repoC.Year, repoC.Color, repoC.InteriorColor,
		repoC.Mileage, repoC.Price, repoC.Status, repoC.Options, repoC.UpdatedAt, repoC.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		} else if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		} else if sqlutil.IsForeignKeyViolationSQL(err) {
			return nil, errors.ErrForeignKey
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
	return r.GetCarByID(ctx, c.ID)
}
