package service

import (
	"github.com/magmaheat/cache-service/internal/repo"
	"github.com/magmaheat/cache-service/pkg/hasher"
	"time"
)

type Services struct {
	Auth  Auth
	Cache Cache
}

type ServicesDependencies struct {
	Repos *repo.Repositories

	AdminToken string
	Hasher     hasher.HashManager
	SignKey    string
	TokenTTL   time.Duration
}

func New(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Repos.Auth, deps.AdminToken, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Cache: NewCacheService(deps.Repos.Cache),
	}
}
