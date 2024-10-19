package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	App
	HTTP
	PG
	RD
	Log
	JWT
}

type App struct {
	Name    string `env-required:"true" env:"APP_NAME"`
	Version string `env-required:"true" env:"APP_VERSION"`
}

type Log struct {
	Level string `env-required:"true" env:"LOG_LEVEL"`
}

type HTTP struct {
	Host string `env-required:"true" env:"HTTP_HOST"`
	Port string `env-required:"true" env:"HTTP_PORT"`
}

type PG struct {
	URL         string `env-required:"true" env:"PG_URL"`
	MaxPoolSize int    `env-required:"true" env:"PG_MAX_POOL_SIZE"`
}

type RD struct {
	URL string `env-required:"true" env:"RD_URL"`
}

type JWT struct {
	SignKey  string        `env-required:"true" env:"JWT_SIGN_KEY"`
	TokenTTL time.Duration `env-required:"true" env:"JWT_TOKEN_TTL"`
}

func New() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("error setup env: %v", err)
	}

	return &cfg
}
