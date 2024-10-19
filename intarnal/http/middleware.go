package http

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/cache-service/intarnal/service"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (m *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
