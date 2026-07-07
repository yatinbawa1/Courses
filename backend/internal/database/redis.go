package database

import (
	"context"

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
