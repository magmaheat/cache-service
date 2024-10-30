package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/magmaheat/cache-service/internal/repo/repoerrs"
	"github.com/magmaheat/cache-service/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (a *AuthRepo) CreateUser(ctx context.Context, login, password string) (string, error) {
	const fn = "repo - pgdb - auth - CreateUser"

	sql, args, _ := a.Builder.
		Insert("users").
		Columns("login, password").
		Values(login, password).
		Suffix("RETURNING login").
		ToSql()

	var userLogin string
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&userLogin)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", repoerrs.ErrAlreadyExists
			}
		}
		log.Errorf("%s - QueryRow: %v", fn, err)
		return "", fmt.Errorf("%s - QueryRow: %v", fn, err)
	}

	return userLogin, nil
}

func (a *AuthRepo) GetUserPassword(ctx context.Context, login string) (string, error) {
	sql, args, _ := a.Builder.
		Select("password").
		From("users").
		Where("login = ?", login).
		ToSql()

	var hash string

	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		log.Errorf("pgdb - auth - GetUserIdAmdPassowrd: %v", err)
		return "", err
	}

	return hash, nil
}

func (a *AuthRepo) AddTokenInBlackList(ctx context.Context, token string) error {
	sql, args, _ := a.Builder.
		Insert("tokens").
		Columns("token").
		Values(token).
		ToSql()

	_, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo - auth - AddTokenInBlackList - Pool.Query: %v", err)
		return err
	}

	return nil
}

func (a *AuthRepo) CheckTokenInBlackList(ctx context.Context, token string) (int, error) {
	sql, args, _ := a.Builder.
		Select("id").
		From("tokens").
		Where("token = ?", token).
		ToSql()

	var id int

	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		log.Errorf("repo - auth - CheckTokenInBlackList - Pool.QueryRow: %v", err)
		return 0, err
	}

	return id, nil
}
