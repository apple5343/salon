package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"
)

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
