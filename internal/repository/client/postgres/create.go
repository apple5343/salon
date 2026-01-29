package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *clientRepository) Create(ctx context.Context, c *service.Client) (*service.Client, error) {
	id := uuid.New().String()
	c.ID = id
	repoC := models.ClientToDatabase(c)
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO clients
		(id, full_name, phone, email, password_hash, passport_series, passport_number, passport_issued_by, birth_date, created_at, updated_at)
		VALUES (:id, :full_name, :phone, :email, :password_hash, :passport_series, :passport_number, :passport_issued_by, :birth_date, :created_at, :updated_at)`,
		repoC)

	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	return r.GetByID(ctx, id)
}
