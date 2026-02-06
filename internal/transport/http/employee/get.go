package employee

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		employee, err := h.service.GetByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EmployeeToHttp(employee))
	}
}

func (h *Handler) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		employee, err := h.service.Profile(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EmployeeToHttp(employee))
	}
}
