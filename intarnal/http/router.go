package http

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/cache-service/intarnal/service"
)

func Init(services *service.Services) *echo.Echo {
	_ = services

	e := echo.New()

	return e
}
