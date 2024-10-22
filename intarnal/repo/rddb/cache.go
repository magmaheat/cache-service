package rddb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/magmaheat/cache-service/intarnal/entity"
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

func (c *CacheRepo) GetDocuments(ctx context.Context) ([]entity.Meta, error) {
	const fn = "repo - cache - GetDocuments"

	keys, err := c.Client.Keys(ctx, "document:*").Result()
	if err != nil {
		log.Errorf("%s - Client.Keys: %v", fn, err)
		return nil, err
	}

	var metaList []entity.Meta

	for _, key := range keys {
		var meta entity.Meta

		result, err := c.Client.HGet(ctx, key, "meta").Result()
		if err != nil {
			log.Errorf("%s - HGet: %v", fn, err)
			return nil, err
		}

		err = json.Unmarshal([]byte(result), &meta)
		if err != nil {
			log.Errorf("%s - json.Unmurshal: %v", fn, err)
			return nil, err
		}

		metaList = append(metaList, meta)
	}

	return metaList, nil
}

func (c *CacheRepo) DeleteDocument(ctx context.Context, key string) (int, error) {
	count, err := c.Client.Del(ctx, key).Result()
	if err != nil {
		log.Errorf("repo - cache - DeleteDocument: %v", err)
		return 0, err
	}

	return int(count), nil
}
