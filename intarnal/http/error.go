package http

import "github.com/labstack/echo/v4"

func newErrorResponse(c echo.Context, errStatus int, message string) {
	report := echo.NewHTTPError(errStatus, message)
	_ = c.JSON(errStatus, report)
}
