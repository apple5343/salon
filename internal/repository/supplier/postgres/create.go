package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *supplierRepository) Create(ctx context.Context, s *service.Supplier) (*service.Supplier, error) {
	repoS := models.SupplierToDatabase(s)
	id := uuid.New().String()
	repoS.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO suppliers
		(id, name, country_code, created_at, updated_at)
		VALUES (:id, :name, :country_code, :created_at, :updated_at)`,
		repoS)

	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	return r.GetByID(ctx, id)
}
