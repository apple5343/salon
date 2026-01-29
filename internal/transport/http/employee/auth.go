package employee

import (
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login() echo.HandlerFunc {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		var req reqBody
		if err := c.Bind(&req); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		accessToken, refreshToken, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func (h *Handler) GetRefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		refreshToken, err := h.service.GetRefreshToken(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"refresh_token": refreshToken,
		})
	}
}

func (h *Handler) GetAccessToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := h.service.GetAccessToken(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"access_token": accessToken,
		})
	}
}
