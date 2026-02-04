package sale

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var s models.Sale
		if err := c.Bind(&s); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceSale, err := models.SaleToService(&s)
		if err != nil {
			return err
		}
		serviceSale, err = h.service.Create(c.Request().Context(), serviceSale)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.SaleToHttp(serviceSale))
	}
}
