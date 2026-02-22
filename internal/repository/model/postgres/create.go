package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *modelRepository) Create(ctx context.Context, m *service.Model) (*service.Model, error) {
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
	return r.GetByID(ctx, id)
}
