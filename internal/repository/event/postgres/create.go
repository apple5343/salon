package postgres

import (
	"context"
	service "salon/internal/models"
	"salon/internal/repository/errors"
	"salon/internal/repository/models"
	sqlutil "salon/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *eventRepository) Create(ctx context.Context, e *service.Event) (*service.Event, error) {
	id := uuid.New().String()
	repoE := models.EventToDatabase(e)
	repoE.ID = id
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO events
		(id, event_type, entity_type, entity_id, actor_id, actor_role, payload, created_at, context)
		VALUES (:id, :event_type, :entity_type, :entity_id, :actor_id, :actor_role, :payload, :created_at, :context)`,
		repoE)

	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, errors.ErrAlreadyExists
		}
		return nil, err
	}
	return r.GetByID(ctx, id)
}
