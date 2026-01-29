package supplier

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var s models.Supplier
		if err := c.Bind(&s); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceS := models.SupplierToService(&s)
		serviceS, err := h.service.Create(c.Request().Context(), serviceS)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.SupplierInternalToHttp(serviceS))
	}
}
