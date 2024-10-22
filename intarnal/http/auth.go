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
	g.DELETE("/auth/:token", r.deleteToken)
}

type registerInput struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"pswd" validate:"required,password"`
	Token    string `json:"token"`
}

func (a *authRoutes) register(c echo.Context) error {
	var input registerInput

	if err := c.Bind(&input); err != nil || input.Token == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return err
	}

	if !a.authRoutes.CheckAdminToken(input.Token) {
		newErrorResponse(c, http.StatusForbidden, "invalid token")
		return fmt.Errorf("invalid token")
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	login, err := a.authRoutes.CreateUser(c.Request().Context(), input.Login, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
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
	Login    string `json:"login" required:"true"`
	Password string `json:"pswd" required:"true"`
}

func (a *authRoutes) auth(c echo.Context) error {
	var input authInput

	if err := c.Bind(&input); err != nil {
		log.Errorf("http - auth - auth - c.Bind: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid body request")
		return err
	}

	token, err := a.authRoutes.GenerateToken(c.Request().Context(), input.Login, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
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

func (a *authRoutes) deleteToken(c echo.Context) error {
	token := c.Param("token")

	if token == "" {
		log.Errorf("http - auth - deleteToken - c.Param: param token empty")
		newErrorResponse(c, http.StatusBadRequest, "param token empty")
		return fmt.Errorf("param token empty")
	}

	err := a.authRoutes.AddTokenInBlackList(c.Request().Context(), token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Response: map[string]bool{
			token: true,
		},
	})
}
