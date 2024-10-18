package pgdb

import "github.com/magmaheat/cache-service/pkg/postgres"

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}
