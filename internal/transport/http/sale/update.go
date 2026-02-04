package sale

import (
	"net/http"
	"salon/internal/transport/http/models"

	service "salon/internal/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Update() echo.HandlerFunc {
	var reqBody struct {
		Status string `json:"status"`
	}
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is empty", errorx.BadRequest)
		}
		if err := c.Bind(&reqBody); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		status, ok := models.SaleStatusMap[reqBody.Status]
		if !ok {
			return errorx.NewError("invalid status", errorx.BadRequest)
		}
		switch status {
		case service.SaleStatusCanceled:
			err := h.service.Cancel(c.Request().Context(), id)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, nil)
		case service.SaleStatusCompleted:
			sale, err := h.service.Complete(c.Request().Context(), id)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, models.SaleToHttp(sale))
		}
		return errorx.NewError("invalid status", errorx.BadRequest)
	}
}
