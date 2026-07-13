package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"salon/pkg/logger"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if nil == err {
				return nil
			}

			if c.Response().Committed {
				return err
			}

			fields := []zap.Field{zap.String("transport", "http"), zap.String("path", c.Path()),
				zap.String("method", c.Request().Method)}
			defer func(err error) {
				logger.FromContextOrDefault(c.Request().Context()).Error(c.Request().Context(), err.Error(), fields...)
			}(err)

			httpErr := &echo.HTTPError{}
			if errors.As(err, &httpErr) {
				fields = append(fields, zap.Int("http_code", httpErr.Code))
				return c.JSON(httpErr.Code, map[string]string{
					"status":  "error",
					"message": fmt.Sprintf("%v", httpErr.Message),
				})
			}

			if commonErr, ok := errorx.ToCommonError(err); ok {
				fields = append(fields, zap.Any("error_code", commonErr.Code()))
				if commonErr.Code() == errorx.Internal {
					fields = append(fields, zap.Int("http_code", http.StatusInternalServerError))
					return c.JSON(
						http.StatusInternalServerError,
						map[string]string{"status": "error", "message": "Internal server error"},
					)
				}

				httpCode := parseCommonCode(commonErr.Code())
				fields = append(fields, zap.Int("http_code", httpCode))
				return c.JSON(
					httpCode,
					map[string]string{"status": "error", "message": commonErr.Error()},
				)
			}

			fields = append(fields, zap.Int("http_code", http.StatusInternalServerError))
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"status": "error", "message": "Internal server error"},
			)
		}
	}
}

func parseCommonCode(code errorx.Code) int {
	switch code {
	case errorx.Unauthorized:
		return http.StatusUnauthorized
	case errorx.Forbidden:
		return http.StatusForbidden
	case errorx.Conflict:
		return http.StatusConflict
	case errorx.OK:
		return http.StatusOK
	case errorx.BadRequest:
		return http.StatusBadRequest
	case errorx.Internal:
		return http.StatusInternalServerError
	case errorx.NotFound:
		return http.StatusNotFound
	case errorx.Cancelled:
		return http.StatusGone
	case errorx.DeadlineExceeded:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
