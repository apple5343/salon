package car

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateBrand() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		var b models.Brand
		if err := c.Bind(&b); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		b.ID = id
		serviceB := models.BrandToService(&b)
		serviceB, err := h.service.UpdateBrand(c.Request().Context(), serviceB)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.BrandInternalToHttp(serviceB))
	}
}

func (h *Handler) UpdateModel() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		var m models.Model
		if err := c.Bind(&m); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		m.ID = id
		serviceM, err := models.ModelToService(&m)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceM, b, err := h.service.UpdateModel(c.Request().Context(), serviceM)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.ModelInternalToHttp(serviceM, b))
	}
}

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
