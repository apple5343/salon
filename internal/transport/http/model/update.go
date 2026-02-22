package model

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
		var m models.Model
		if err := c.Bind(&m); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		m.ID = id
		serviceM, err := models.ModelToService(&m)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceM, b, err := h.service.Update(c.Request().Context(), serviceM)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.ModelInternalToHttp(serviceM, b))
	}
}
