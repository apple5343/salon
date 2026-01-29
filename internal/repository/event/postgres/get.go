package postgres

import (
	"context"
	"database/sql"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
)

func (r *eventRepository) GetByID(ctx context.Context, id string) (*service.Event, error) {
	var e models.Event
	err := r.db.GetContext(ctx, &e, `SELECT * FROM events WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return models.EventToService(&e), nil
}
