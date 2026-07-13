package middlewares

import (
	"net/http"
	"net/http/httptest"
	"salon/pkg/logger"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	l := logger.NewLogger(&logger.Config{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	LoggerMiddleware(l)(func(c echo.Context) error {
		_, ok := logger.FromContext(c.Request().Context())
		require.True(t, ok)
		return nil
	})(c)
}
