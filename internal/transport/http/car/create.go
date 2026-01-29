package car

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateBrand() echo.HandlerFunc {
	return func(c echo.Context) error {
		var b models.Brand
		if err := c.Bind(&b); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceB := models.BrandToService(&b)
		serviceB, err := h.service.CreateBrand(c.Request().Context(), serviceB)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.BrandInternalToHttp(serviceB))
	}
}

func (h *Handler) CreateModel() echo.HandlerFunc {
	return func(c echo.Context) error {
		var m models.Model
		if err := c.Bind(&m); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceM, err := models.ModelToService(&m)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceM, b, err := h.service.CreateModel(c.Request().Context(), serviceM)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.ModelInternalToHttp(serviceM, b))
	}
}

func (h *Handler) CreateCar() echo.HandlerFunc {
	return func(c echo.Context) error {
		var car models.Car
		if err := c.Bind(&car); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceCar, err := models.CarToService(&car)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceCar, model, brand, supplier, err := h.service.CreateCar(c.Request().Context(), serviceCar)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.CarInternalToHttp(serviceCar, model, brand, supplier))
	}
}
