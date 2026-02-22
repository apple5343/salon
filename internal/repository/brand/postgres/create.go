package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *brandRepository) Create(ctx context.Context, b *service.Brand) (*service.Brand, error) {
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
	return r.GetByID(ctx, id)
}
