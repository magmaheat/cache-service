package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magmaheat/cache-service/intarnal/service"
	"os"

	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Response struct {
	Errors   *ErrorResponse `json:"error,omitempty"`
	Response interface{}    `json:"response,omitempty"`
	Data     interface{}    `json:"data,omitempty"`
}

func Init(services *service.Services) *echo.Echo {
	handler := echo.New()

	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := handler.Group("/api")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	h := handler.Group("/api", authMiddleware.UserIdentity)
	{
		_ = h
	}

	//h.Any("/api/*", func(c echo.Context) error {
	//	newErrorResponse(c, http.StatusMethodNotAllowed, "Method not allowed")
	//	return nil
	//})

	return handler
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("http - router - setLogsFile: %v", err)
	}

	return file
}
