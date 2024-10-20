package rddb

import (
	"context"
	"fmt"
	"github.com/magmaheat/cache-service/pkg/redis"
	log "github.com/sirupsen/logrus"
)

type CacheRepo struct {
	*redis.Redis
}

func NewCacheRepo(rd *redis.Redis) *CacheRepo {
	return &CacheRepo{rd}
}

func (c *CacheRepo) SaveJson(ctx context.Context, key, field, data string) error {
	_, err := c.Client.HSet(ctx, key, field, data).Result()
	if err != nil {
		log.Errorf("repo - cache - SaveJson - HSet: %v", err)
		return fmt.Errorf("error record meta: %s", err)
	}

	return nil
}

func (c *CacheRepo) SaveFile(ctx context.Context, key, field string, file []byte) error {
	_, err := c.Client.HSet(ctx, key, field, file).Result()
	if err != nil {
		log.Errorf("repo - cache - SaveFile - HSet: %v", err)
		return fmt.Errorf("error record file: %s", err)
	}

	return nil
}

func (c *CacheRepo) CheckId(ctx context.Context, id string) (bool, error) {
	exists, err := c.Client.Exists(ctx, id).Result()
	if err != nil {
		log.Errorf("repo - cache - CheckId - Client.Exists: %v", err)
		return false, err
	}

	return exists == 1, nil
}

func (c *CacheRepo) GetDocument(ctx context.Context, id string) (map[string]string, error) {
	allFields, err := c.Client.HGetAll(ctx, id).Result()
	if err != nil {
		log.Errorf("repo - cache - GetDocument - HGetAll: %v", err)
		return map[string]string{}, err
	}

	return allFields, nil
}
