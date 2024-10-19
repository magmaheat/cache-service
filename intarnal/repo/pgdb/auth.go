package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/magmaheat/cache-service/intarnal/repo/repoerrs"
	"github.com/magmaheat/cache-service/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (a *AuthRepo) CreateUser(ctx context.Context, username, password string) (int, error) {
	const fn = "repo - pgdb - auth - CreateUser"

	sql, args, _ := a.Builder.
		Insert("users").
		Columns("username, password").
		Values(username, password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		log.Errorf("%s - QueryRow: %v", fn, err)
		return 0, fmt.Errorf("%s - QueryRow: %v", fn, err)
	}

	return id, nil
}

func (a *AuthRepo) GetUserIdAndPassword(ctx context.Context, username, password string) (int, string, error) {
	sql, args, _ := a.Builder.
		Select("id, password").
		From("users").
		Where("username = ?", username).
		ToSql()

	var hash string
	var id int

	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", repoerrs.ErrNotFound
		}
		log.Errorf("pgdb - auth - GetUserIdAmdPassowrd: %v", err)
		return 0, "", err
	}

	return id, hash, nil
}
