package car

import (
	"net/http"
	"salon/internal/transport/http/models"
	ctxutil "salon/internal/utils/context"

	service "salon/internal/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetCarByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		car, model, brand, supplier, err := h.service.GetCarByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			return c.JSON(http.StatusOK, models.CarInternalToHttp(car, model, brand, supplier))
		} else {
			return c.JSON(http.StatusOK, models.CarPublicToHttp(car, model, brand, supplier))
		}
	}
}
