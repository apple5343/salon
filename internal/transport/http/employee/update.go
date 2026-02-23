package employee

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		var e models.Employee
		if err := c.Bind(&e); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		e.ID = id
		serviceE, err := models.EmployeeToService(&e)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceE, err = h.service.Update(c.Request().Context(), serviceE)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EmployeeToHttp(serviceE))
	}
}

func (h *Handler) Hire() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		e, err := h.service.Hire(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EmployeeToHttp(e))
	}
}
