package pgdb

import (
	"github.com/magmaheat/cache-service/pkg/postgres"
)

type FilesRepo struct {
	*postgres.Postgres
}

func NewFilesRepo(pg *postgres.Postgres) *FilesRepo {
	return &FilesRepo{pg}
}
