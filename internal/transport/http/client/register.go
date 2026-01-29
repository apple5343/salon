package client

import (
	"net/http"
	"salon/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var client models.Client
		if err := c.Bind(&client); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceClient, err := models.ClientToService(&client)
		if err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		serviceClient, err = h.service.Register(c.Request().Context(), serviceClient)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.ClientToHttp(serviceClient))
	}
}
