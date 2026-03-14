package car

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateCar() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		var car models.Car
		if err := c.Bind(&car); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		car.ID = id
		serviceCar, err := models.CarToService(&car)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceCar, model, brand, supplier, err := h.service.UpdateCar(c.Request().Context(), serviceCar)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.CarInternalToHttp(serviceCar, model, brand, supplier))
	}
}
