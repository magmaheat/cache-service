package http

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/magmaheat/cache-service/intarnal/service"
	"net/http"
)

type authRoutes struct {
	authRoutes service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authRoutes: authService,
	}

	g.POST("/register", r.register)
	g.POST("/auth", r.auth)
}

type registerInput struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
	Token    string `json:"token"`
}

func (a *authRoutes) register(c echo.Context) error {
	var input registerInput

	if err := c.Bind(&input); err != nil || input.Token == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid tut body request")
		return err
	}

	if !a.authRoutes.CheckAdminToken(input.Token) {
		newErrorResponse(c, http.StatusForbidden, "invalid token")
		return fmt.Errorf("invalid token")
	}

	//TODO add validate

	login, err := a.authRoutes.CreateUser(c.Request().Context(), input.Login, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Login string `json:"login"`
	}

	return c.JSON(http.StatusOK, Response{
		Response: response{
			Login: login,
		},
	})
}

type authInput struct {
	Username string `json:"login"`
	Password string `json:"pswd"`
}

func (a *authRoutes) auth(c echo.Context) error {
	var input authInput

	if err := c.Bind(&input); err != nil {
		log.Errorf("http - auth - auth - c.Bind: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return err
	}

	token, err := a.authRoutes.GenerateToken(c.Request().Context(), input.Username, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}

		if errors.Is(err, service.ErrInvalidPassword) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}

		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, Response{
		Response: response{
			Token: token,
		},
	})
}
