package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/magmaheat/cache-service/intarnal/repo"
	"github.com/magmaheat/cache-service/intarnal/repo/repoerrs"
	"github.com/magmaheat/cache-service/pkg/hasher"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type Auth interface {
	CreateUser(ctx context.Context, username, password string) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
}

type AuthService struct {
	authRepo repo.Auth

	hasher   hasher.HashManager
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(authRepo repo.Auth, hasher hasher.HashManager, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		hasher:   hasher,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (a *AuthService) CreateUser(ctx context.Context, username, password string) (int, error) {
	hash, err := a.hasher.HashPassword(password)
	if err != nil {
		log.Errorf("service - auth - CreateUser - HashPassword: %v", err)
		return 0, err
	}

	id, err := a.authRepo.CreateUser(ctx, username, hash)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}
		return 0, err
	}

	return id, nil
}

func (a *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	id, hash, err := a.authRepo.GetUserIdAndPassword(ctx, username, password)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	if !a.hasher.CheckPassword(hash, password) {
		return "", ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: id,
	})

	tokenString, err := token.SignedString([]byte(a.signKey))
	if err != nil {
		log.Errorf("service - auth - GenerateToken - SignedString: %v", err)
		return "", err
	}

	return tokenString, nil
}
