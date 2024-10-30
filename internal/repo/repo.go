package repo

import (
	"context"
	"github.com/magmaheat/cache-service/internal/entity"
	"github.com/magmaheat/cache-service/internal/repo/pgdb"
	"github.com/magmaheat/cache-service/internal/repo/rddb"
	"github.com/magmaheat/cache-service/pkg/postgres"
	"github.com/magmaheat/cache-service/pkg/redis"
)

type Auth interface {
	CreateUser(ctx context.Context, login, password string) (string, error)
	GetUserPassword(ctx context.Context, login string) (string, error)
	AddTokenInBlackList(ctx context.Context, token string) error
	CheckTokenInBlackList(ctx context.Context, token string) (int, error)
}

type Cache interface {
	SaveJson(ctx context.Context, key, field, data string) error
	SaveFile(ctx context.Context, key, field string, file []byte) error
	CheckId(ctx context.Context, id string) (bool, error)
	GetDocument(ctx context.Context, id string) (map[string]string, error)
	GetDocuments(ctx context.Context) ([]entity.Meta, error)
	DeleteDocument(ctx context.Context, key string) (int, error)
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
