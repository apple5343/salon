package employee

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var e models.Employee
		if err := c.Bind(&e); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceE, err := models.EmployeeToService(&e)
		if err != nil {
			return err
		}
		serviceE, err = h.service.Register(c.Request().Context(), serviceE)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.EmployeeToHttp(serviceE))
	}
}
