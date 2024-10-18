package rddb

import "github.com/magmaheat/cache-service/pkg/redis"

type CacheRepo struct {
	*redis.Redis
}

func NewCacheRepo(rd *redis.Redis) *CacheRepo {
	return &CacheRepo{rd}
}
