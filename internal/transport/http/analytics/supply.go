package analytics

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Supply() echo.HandlerFunc {
	return func(c echo.Context) error {
		dateFrom, dateTo, err := timeRangeFromRequest(c)
		if err != nil {
			return ErrInvalidTimeFormat
		}
		supply, err := h.service.Supply(c.Request().Context(), dateFrom, dateTo)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.SupplyAnalyticsToHttp(supply))
	}
}
