package event

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEvents() echo.HandlerFunc {
	return func(c echo.Context) error {
		filters, err := eventFiltersFromRequest(c)
		if err != nil {
			return err
		}
		events, err := h.service.GetEvents(c.Request().Context(), filters)
		if err != nil {
			return err
		}
		result := make([]*models.Event, len(events))
		for i, e := range events {
			result[i] = models.EventToHttp(e)
		}
		c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(events)))
		return c.JSON(http.StatusOK, result)
	}
}

func eventFiltersFromRequest(c echo.Context) (*service.EventFilters, error) {
	filters := &service.EventFilters{}
	if eventType := c.QueryParam("event_type"); eventType != "" {
		t, ok := models.EventTypeMap[eventType]
		if !ok {
			return nil, errorx.NewError("invalid event type", errorx.BadRequest)
		}
		filters.Type = &t
	}
	if entityType := c.QueryParam("entity_type"); entityType != "" {
		t, ok := models.EntityTypeMap[entityType]
		if !ok {
			return nil, errorx.NewError("invalid entity type", errorx.BadRequest)
		}
		filters.EntityType = &t
	}
	if entityID := c.QueryParam("entity_id"); entityID != "" {
		filters.EntityID = &entityID
	}
	if actorID := c.QueryParam("actor_id"); actorID != "" {
		filters.ActorID = &actorID
	}
	if actorRole := c.QueryParam("actor_role"); actorRole != "" {
		t, ok := models.EmployeeRoleMap[actorRole]
		if !ok {
			return nil, errorx.NewError("invalid actor role", errorx.BadRequest)
		}
		filters.ActorRole = &t
	}
	if limit := c.QueryParam("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Limit = &limitInt
	}
	if offset := c.QueryParam("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, errorx.NewError(err.Error(), errorx.BadRequest)
		}
		filters.Offset = &offsetInt
	}
	return filters, nil
}
