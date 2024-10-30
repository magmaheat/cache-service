package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magmaheat/cache-service/internal/service"
	"github.com/magmaheat/cache-service/pkg/validator"
	"net/http"
	"os"

	_ "github.com/magmaheat/cache-service/docs"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Response struct {
	Errors   ErrorResponse `json:"error,omitempty"`
	Response interface{}   `json:"response,omitempty"`
	Data     interface{}   `json:"data,omitempty"`
}

func Init(services *service.Services) *echo.Echo {
	handler := echo.New()

	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.Validator = validator.NewCustomValidator()

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := handler.Group("/api")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	fileGroup := handler.Group("/api", authMiddleware.UserIdentity)
	{
		NewFilesRouter(fileGroup, services)
	}

	handler.Any("/api/*", func(c echo.Context) error {
		log.Debugf("method %s not implementation for URL %s", c.Request().Method, c.Request().URL)
		newErrorResponse(c, http.StatusNotImplemented, "Method not allowed")
		return nil
	})

	return handler
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("http - router - setLogsFile: %v", err)
	}

	return file
}
