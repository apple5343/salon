package models

import (
	"salon/internal/models"
	"time"
)

var EventTypeMap = map[string]models.EventType{
	"created": models.EventTypeCreated,
	"updated": models.EventTypeUpdated,
	"deleted": models.EventTypeDeleted,
}

var EntityTypeMap = map[string]models.EntityType{
	"car":      models.EntityTypeCar,
	"supplier": models.EntityTypeSupplier,
	"model":    models.EntityTypeModel,
	"brand":    models.EntityTypeBrand,
	"employee": models.EntityTypeEmployee,
}

var EmployeeRoleMap = map[string]models.EmployeeRole{
	"admin": models.EmployeeRoleAdmin,
	"user":  models.EmployeeRoleManager,
}

type Event struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"event_type"`
	EntityType string                 `json:"entity_type"`
	EntityID   string                 `json:"entity_id"`
	ActorID    string                 `json:"actor_id"`
	ActorRole  string                 `json:"actor_role"`
	Payload    map[string]interface{} `json:"payload"`
	CreatedAt  time.Time              `json:"created_at"`
	Context    map[string]interface{} `json:"context"`
}

func EventToHttp(e *models.Event) *Event {
	return &Event{
		ID:         e.ID,
		Type:       string(e.Type),
		EntityType: string(e.EntityType),
		EntityID:   e.EntityID,
		ActorID:    e.ActorID,
		ActorRole:  e.ActorRole,
		Payload:    e.Payload,
		CreatedAt:  e.CreatedAt,
		Context:    e.Context,
	}
}
