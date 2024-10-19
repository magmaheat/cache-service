package repo

import (
	"context"
	"github.com/magmaheat/cache-service/intarnal/repo/pgdb"
	"github.com/magmaheat/cache-service/intarnal/repo/rddb"
	"github.com/magmaheat/cache-service/pkg/postgres"
	"github.com/magmaheat/cache-service/pkg/redis"
)

type Auth interface {
	CreateUser(ctx context.Context, username, password string) (int, error)
	GetUserIdAndPassword(ctx context.Context, username, password string) (int, string, error)
}

type Cache interface {
}

type Files interface {
}

type Repositories struct {
	Auth
	Files
	Cache
}

func New(pg *postgres.Postgres, rd *redis.Redis) *Repositories {
	return &Repositories{
		Auth:  pgdb.NewAuthRepo(pg),
		Files: pgdb.NewFilesRepo(pg),
		Cache: rddb.NewCacheRepo(rd),
	}
}
