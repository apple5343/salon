package analytics

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Sales() echo.HandlerFunc {
	return func(c echo.Context) error {
		dateFrom, dateTo, err := timeRangeFromRequest(c)
		if err != nil {
			return err
		}
		sales, err := h.service.Sales(c.Request().Context(), dateFrom, dateTo)
		if err != nil {
			return ErrInvalidTimeFormat
		}
		return c.JSON(http.StatusOK, sales)
	}
}
