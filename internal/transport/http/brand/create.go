package brand

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var b models.Brand
		if err := c.Bind(&b); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceB := models.BrandToService(&b)
		serviceB, err := h.service.Create(c.Request().Context(), serviceB)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.BrandInternalToHttp(serviceB))
	}
}
