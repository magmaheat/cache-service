package redis

import "time"

type Options func(r *Redis)

func connAttempts(attempts int) Options {
	return func(r *Redis) {
		r.connAttempts = attempts
	}
}

func connTimeout(timeout time.Duration) Options {
	return func(r *Redis) {
		r.connTimeout = timeout
	}
}
