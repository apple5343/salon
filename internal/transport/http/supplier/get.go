package supplier

import (
	"net/http"
	service "salon/internal/models"
	"salon/internal/transport/http/models"
	ctxutil "salon/internal/utils/context"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		supplier, err := h.service.GetByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		if role := ctxutil.UserRoleFromContext(c.Request().Context()); role == string(service.EmployeeRoleAdmin) || role == string(service.EmployeeRoleManager) {
			return c.JSON(http.StatusOK, models.SupplierInternalToHttp(supplier))
		}
		return c.JSON(http.StatusOK, models.SupplierPublicToHttp(supplier))
	}
}
