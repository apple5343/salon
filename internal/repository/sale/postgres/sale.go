package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *saleRepository) Create(ctx context.Context, s *service.Sale) (*service.Sale, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var carStatus string
	err = tx.QueryRowContext(ctx, `SELECT status FROM cars WHERE id = $1 FOR UPDATE`, s.CarID).Scan(&carStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	if carStatus != string(service.CarStatusAvailable) {
		return nil, errors.ErrCarNotAvailable
	}

	id := uuid.New().String()
	repoS := models.SaleToDatabase(s)
	repoS.ID = id
	_, err = tx.NamedExecContext(ctx, `INSERT INTO sales
		(id, car_id, client_id, employee_id, original_price, discount_amount, discount_percent,
		payment_type, status, notes, created_at, updated_at) 
		VALUES (:id, :car_id, :client_id, :employee_id, :original_price, :discount_amount, :discount_percent, 
		:payment_type, :status, :notes, :created_at, :updated_at)`,
		repoS)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		} else if sqlutil.IsForeignKeyViolationSQL(err) {
			return nil, errors.ErrForeignKey
		}
		return nil, err
	}
	_, err = tx.ExecContext(ctx, `UPDATE cars SET status = $1 WHERE id = $2`, string(service.CarStatusBooked), s.CarID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *saleRepository) Complete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := r.clock.Now()
	res := tx.QueryRowContext(ctx, "UPDATE sales SET status = $1, updated_at = $2, sale_date = $3 WHERE id = $4 AND status = $5 RETURNING car_id",
		string(service.SaleStatusCompleted), now, now, id, string(service.SaleStatusPending))

	var carID string
	if err = res.Scan(&carID); err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE cars SET status = $1 WHERE id = $2", string(service.CarStatusSold), carID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *saleRepository) Cancel(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := r.clock.Now()
	res := tx.QueryRowContext(ctx, "UPDATE sales SET status = $1 , updated_at = $2 WHERE id = $3 AND status = $4 RETURNING car_id",
		string(service.SaleStatusCanceled), now, id, string(service.SaleStatusPending))

	var carID string
	if err = res.Scan(&carID); err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE cars SET status = $1 WHERE id = $2", string(service.CarStatusAvailable), carID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
