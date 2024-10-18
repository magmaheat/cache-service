package service

import "github.com/magmaheat/cache-service/intarnal/repo"

type Cache interface {
}

type CacheService struct {
	cacheRepo repo.Cache
}

func NewCacheService(cacheRepo repo.Cache) *CacheService {
	return &CacheService{cacheRepo}
}
