package middlewares

import (
	"salon/internal/config"
	ctxutil "salon/internal/utils/context"
	"salon/pkg/jwt"
	"strings"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

const (
	bearerPrefix = "Bearer "
)

func AuthMiddleware(jwtConfig *config.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errorx.NewError("token is empty", errorx.Unauthorized)
			}
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			userClaims, err := jwt.VerifyToken(token, []byte(jwtConfig.AccessSecret))
			if err != nil {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			ctx := ctxutil.ContextWithUserID(c.Request().Context(), userClaims.ID)
			ctx = ctxutil.ContextWithUserRole(ctx, userClaims.Role)
			ctx = ctxutil.ContextWithUserToken(ctx, token)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func SoftAuthMiddleware(jwtConfig *config.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return next(c)
			}
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			userClaims, err := jwt.VerifyToken(token, []byte(jwtConfig.AccessSecret))
			if err != nil {
				return next(c)
			}
			ctx := ctxutil.ContextWithUserID(c.Request().Context(), userClaims.ID)
			ctx = ctxutil.ContextWithUserRole(ctx, userClaims.Role)
			ctx = ctxutil.ContextWithUserToken(ctx, token)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func TokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errorx.NewError("token is empty", errorx.Unauthorized)
			}
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			ctx := ctxutil.ContextWithUserToken(c.Request().Context(), token)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
