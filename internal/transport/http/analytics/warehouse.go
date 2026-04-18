package analytics

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Warehouse() echo.HandlerFunc {
	return func(c echo.Context) error {
		dateFrom, dateTo, err := timeRangeFromRequest(c)
		if err != nil {
			return ErrInvalidTimeFormat
		}
		warehouse, err := h.service.Warehouse(c.Request().Context(), dateFrom, dateTo)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, warehouse)
	}
}
