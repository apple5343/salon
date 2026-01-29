package models

import "time"

type EventType string

var (
	EventTypeCreated EventType = "created"
	EventTypeUpdated EventType = "updated"
	EventTypeDeleted EventType = "deleted"
)

type EntityType string

var (
	EntityTypeBrand    EntityType = "brand"
	EntityTypeModel    EntityType = "model"
	EntityTypeCar      EntityType = "car"
	EntityTypeEmployee EntityType = "employee"
	EntityTypeSupplier EntityType = "supplier"
)

type Event struct {
	ID         string
	Type       EventType
	EntityType EntityType
	EntityID   string
	ActorID    string
	ActorRole  string
	Payload    map[string]interface{}
	CreatedAt  time.Time
	Context    map[string]interface{}
}

type EventFilters struct {
	Type       *EventType
	EntityType *EntityType
	EntityID   *string `validate:"omitempty,uuid"`
	ActorID    *string `validate:"omitempty,uuid"`
	ActorRole  *EmployeeRole
	BaseList
}

func (f *EventFilters) Validate() error {
	return Validator().Struct(f)
}

func BrandPayload(b *Brand) map[string]interface{} {
	return map[string]interface{}{
		"id":           b.ID,
		"name":         b.Name,
		"country_code": b.CountryCode,
		"description":  b.Description,
		"created_at":   b.CreatedAt,
		"updated_at":   b.UpdatedAt,
	}
}

func ModelPayload(m *Model) map[string]interface{} {
	return map[string]interface{}{
		"id":                        m.ID,
		"brand_id":                  m.BrandID,
		"name":                      m.Name,
		"generation":                m.Generation,
		"body_type":                 m.BodyType,
		"transmission_type":         m.TransmissionType,
		"fuel_type":                 m.FuelType,
		"engine_displacement":       m.EngineDisplacement,
		"power_hp":                  m.PowerHP,
		"drive_type":                m.DriveType,
		"base_price":                m.BasePrice,
		"technical_characteristics": m.TechnicCharacteristics,
		"created_at":                m.CreatedAt,
		"updated_at":                m.UpdatedAt,
	}
}

func EmployeePayload(e *Employee) map[string]interface{} {
	return map[string]interface{}{
		"id":         e.ID,
		"full_name":  e.FullName,
		"phone":      e.Phone,
		"email":      e.Email,
		"role":       e.Role,
		"status":     e.Status,
		"hire_date":  e.HireDate,
		"created_at": e.CreatedAt,
		"updated_at": e.UpdatedAt,
	}
}

func SupplierPayload(s *Supplier) map[string]interface{} {
	return map[string]interface{}{
		"id":           s.ID,
		"name":         s.Name,
		"country_code": s.CountryCode,
		"created_at":   s.CreatedAt,
		"updated_at":   s.UpdatedAt,
	}
}

func CarPayload(c *Car) map[string]interface{} {
	return map[string]interface{}{
		"id":             c.ID,
		"model_id":       c.ModelID,
		"supplier_id":    c.SupplierID,
		"vin":            c.Vin,
		"year":           c.Year,
		"color":          c.Color,
		"interior_color": c.InteriorColor,
		"mileage":        c.Mileage,
		"price":          c.Price,
		"status":         c.Status,
		"options":        c.Options,
		"created_at":     c.CreatedAt,
		"updated_at":     c.UpdatedAt,
	}
}
