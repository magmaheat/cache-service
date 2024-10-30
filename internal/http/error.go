package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
	ErrCannotParseToken  = fmt.Errorf("cannot parse token")
)

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func newErrorResponse(c echo.Context, status int, message string) {
	c.JSON(status, Response{
		Errors: ErrorResponse{
			Code: status,
			Text: message,
		},
	})
}
