package analytics

import (
	"net/http"
	"salon/internal/transport/http/models"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Warehouse() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dateFrom, dateTo *time.Time
		if dateFromStr := c.QueryParam("date_from"); dateFromStr != "" {
			t, err := time.Parse(models.TimeLayout, dateFromStr)
			if err != nil {
				return err
			}
			dateFrom = &t
		}
		if dateToStr := c.QueryParam("date_to"); dateToStr != "" {
			t, err := time.Parse(models.TimeLayout, dateToStr)
			if err != nil {
				return err
			}
			dateTo = &t
		}
		warehouse, err := h.service.Warehouse(c.Request().Context(), dateFrom, dateTo)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, warehouse)
	}
}
