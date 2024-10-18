package service

import (
	"github.com/magmaheat/cache-service/intarnal/repo"
	"time"
)

type Auth interface {
}

type AuthService struct {
	authRepo repo.Auth

	SignKey  string
	TokenTTL time.Duration
}

func NewAuthService(authRepo repo.Auth, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		SignKey:  signKey,
		TokenTTL: tokenTTL,
	}
}
