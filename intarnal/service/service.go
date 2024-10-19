package service

import (
	"github.com/magmaheat/cache-service/intarnal/repo"
	"github.com/magmaheat/cache-service/pkg/hasher"
	"time"
)

type Services struct {
	Auth  Auth
	Files Files
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
		Files: NewFilesService(deps.Repos.Files),
		Cache: NewCacheService(deps.Repos.Cache),
	}
}
