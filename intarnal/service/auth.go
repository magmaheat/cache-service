package service

import (
	"context"
	"errors"
	"fmt"
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
	CheckAdminToken(token string) bool
	CreateUser(ctx context.Context, login, password string) (string, error)
	GenerateToken(ctx context.Context, login, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type AuthService struct {
	authRepo   repo.Auth
	adminToken string

	hasher   hasher.HashManager
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(
	authRepo repo.Auth,
	adminToken string,
	hasher hasher.HashManager,
	signKey string,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		adminToken: adminToken,
		hasher:     hasher,
		signKey:    signKey,
		tokenTTL:   tokenTTL,
	}
}

func (a *AuthService) CheckAdminToken(token string) bool {
	return a.adminToken == token
}

func (a *AuthService) CreateUser(ctx context.Context, login, password string) (string, error) {
	hash, err := a.hasher.HashPassword(password)
	if err != nil {
		log.Errorf("service - auth - CreateUser - HashPassword: %v", err)
		return "", err
	}

	userLogin, err := a.authRepo.CreateUser(ctx, login, hash)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return "", ErrUserAlreadyExists
		}
		return "", err
	}

	return userLogin, nil
}

func (a *AuthService) GenerateToken(ctx context.Context, login, password string) (string, error) {
	id, hash, err := a.authRepo.GetUserIdAndPassword(ctx, login)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrUserNotFound
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

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("cannot parse token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("cannot parse token")
	}

	return claims.UserId, nil
}
