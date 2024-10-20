package repo

import (
	"context"
	"github.com/magmaheat/cache-service/intarnal/repo/pgdb"
	"github.com/magmaheat/cache-service/intarnal/repo/rddb"
	"github.com/magmaheat/cache-service/pkg/postgres"
	"github.com/magmaheat/cache-service/pkg/redis"
)

type Auth interface {
	CreateUser(ctx context.Context, login, password string) (string, error)
	GetUserIdAndPassword(ctx context.Context, login string) (int, string, error)
}

type Cache interface {
	SaveJson(ctx context.Context, key, field, data string) error
	SaveFile(ctx context.Context, key, field string, file []byte) error
	CheckId(ctx context.Context, id string) (bool, error)
	GetDocument(ctx context.Context, id string) (map[string]string, error)
}

type Repositories struct {
	Auth
	Cache
}

func New(pg *postgres.Postgres, rd *redis.Redis) *Repositories {
	return &Repositories{
		Auth:  pgdb.NewAuthRepo(pg),
		Cache: rddb.NewCacheRepo(rd),
	}
}
