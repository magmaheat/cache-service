package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Redis struct {
	connAttempts int
	connTimeout  time.Duration

	Client *redis.Client
}

func New(url string, opts ...Options) *Redis {
	rdb := &Redis{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(rdb)
	}

	optsClient, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("redis - New - ParseURL: %v", err)
	}

	rdb.Client = redis.NewClient(optsClient)
	ctx := context.Background()

	for rdb.connAttempts > 0 {
		_, err = rdb.Client.Ping(ctx).Result()

		if err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts: %d", rdb.connAttempts)
		time.Sleep(rdb.connTimeout)
		rdb.connAttempts--
	}

	if err != nil {
		log.Fatalf("redis - New - rdb.Client.Ping: %v", err)
	}

	return rdb
}
