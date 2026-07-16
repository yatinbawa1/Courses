package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOTPRepo struct {
	rdbs *redis.Client
}

func NewRedisOTPRepo(rdbs *redis.Client) *RedisOTPRepo {
	return &RedisOTPRepo{rdbs}
}

func (r *RedisOTPRepo) SaveOTP(ctx context.Context, email string, code string, ttl time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)

	err := r.rdbs.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save OTP to redis: %w", err)
	}

	return nil
}

func (r *RedisOTPRepo) GetOTPForUser(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	storedCode, err := r.rdbs.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve OTP from redis: %w", err)
	}

	return storedCode, nil
}

func (r *RedisOTPRepo) DeleteOTPForUser(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)
	_, err := r.rdbs.Del(ctx, key).Result()

	if err != nil {
		return err
	}

	return nil
}
