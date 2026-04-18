package analytics

import (
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Supplier() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		dateFrom, dateTo, err := timeRangeFromRequest(c)
		if err != nil {
			return ErrInvalidTimeFormat
		}
		supplier, err := h.service.Supplier(c.Request().Context(), id, dateFrom, dateTo)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, supplier)
	}
}
