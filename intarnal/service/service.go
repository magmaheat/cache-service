package service

import (
	"github.com/magmaheat/cache-service/intarnal/repo"
	"time"
)

type Services struct {
	Auth  Auth
	Files Files
	Cache Cache
}

type ServicesDependencies struct {
	Repos *repo.Repositories

	SignKey  string
	TokenTTL time.Duration
}

func New(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Repos.Auth, deps.SignKey, deps.TokenTTL),
		Files: NewFilesService(deps.Repos.Files),
		Cache: NewCacheService(deps.Repos.Cache),
	}
}
