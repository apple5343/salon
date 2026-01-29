package postgres

import (
	"context"
	"fmt"
	service "salon/internal/models"
	"salon/internal/repository/models"
	"strings"
)

func (r *eventRepository) GetEventsByFilter(ctx context.Context, filter *service.EventFilters) ([]*service.Event, error) {
	var events []*models.Event

	query := `SELECT * FROM events`
	conditions := []string{}
	args := []interface{}{}
	addCondition := func(condition string, val interface{}) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf("%s $%d", condition, len(args)))
	}
	if filter.Type != nil {
		addCondition("event_type =", string(*filter.Type))
	}
	if filter.EntityType != nil {
		addCondition("entity_type =", string(*filter.EntityType))
	}
	if filter.EntityID != nil {
		addCondition("entity_id =", *filter.EntityID)
	}
	if filter.ActorID != nil {
		addCondition("actor_id =", *filter.ActorID)
	}
	if filter.ActorRole != nil {
		addCondition("actor_role =", *filter.ActorRole)
	}
	if len(conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}
	fmt.Println(query)
	if err := r.db.SelectContext(ctx, &events, query, args...); err != nil {
		return nil, err
	}
	result := make([]*service.Event, len(events))
	for i, e := range events {
		result[i] = models.EventToService(e)
	}
	return result, nil
}
