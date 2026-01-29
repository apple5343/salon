package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"salon/internal/models"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

type Event struct {
	ID         string    `json:"id" db:"id"`
	Type       string    `json:"event_type" db:"event_type"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	ActorID    string    `json:"actor_id" db:"actor_id"`
	ActorRole  string    `json:"actor_role" db:"actor_role"`
	Payload    JSONB     `json:"payload" db:"payload"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Context    JSONB     `json:"context" db:"context"`
}

func EventToDatabase(e *models.Event) *Event {
	return &Event{
		ID:         e.ID,
		Type:       string(e.Type),
		EntityType: string(e.EntityType),
		EntityID:   e.EntityID,
		ActorID:    e.ActorID,
		ActorRole:  e.ActorRole,
		Payload:    JSONB(e.Payload),
		CreatedAt:  e.CreatedAt,
		Context:    JSONB(e.Context),
	}
}

func EventToService(e *Event) *models.Event {
	return &models.Event{
		ID:         e.ID,
		Type:       models.EventType(e.Type),
		EntityType: models.EntityType(e.EntityType),
		EntityID:   e.EntityID,
		ActorID:    e.ActorID,
		ActorRole:  e.ActorRole,
		Payload:    e.Payload,
		CreatedAt:  e.CreatedAt,
		Context:    e.Context,
	}
}
