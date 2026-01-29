package event

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		event, err := h.service.GetEventByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EventToHttp(event))
	}
}
