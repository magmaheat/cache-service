package http

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/cache-service/intarnal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (m *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const fn = ""

		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("http - middleware - UserIdentity - bearerToken")
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		id, err := m.authService.ParseToken(token)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return err
		}

		c.Set("userId", id)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
