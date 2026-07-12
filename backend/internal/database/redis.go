package database

import (
	"context"
	"courses/internal/config"
	"errors"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string, db int) (*redis.Client, error) {
	rdbs := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdbs.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdbs, nil
}

func NewRedisOnlineClient() (* redis.Client, error) {
	ErrRedisUrl := errors.New("Unable to find REDIS_URL in env")
	if len(config.REDIS_URL) == 0 {
		return nil, ErrRedisUrl
	}

	opt, err := redis.ParseURL(config.REDIS_URL)
	if err != nil {
		return nil, ErrRedisUrl
	}

	return redis.NewClient(opt), nil
}