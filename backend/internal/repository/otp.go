package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPRepo interface {
	SaveOTP(ctx context.Context, email string, code string, tll time.Duration) error
	VerifyOTP(ctx context.Context, email string, code string) (bool, error)
}

// Interface is Being Implemented
// For Redis, Can also do
// For Any other DBMS by Just Implementing
// Features in OTPRepo

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

func (r *RedisOTPRepo) VerifyOTP(ctx context.Context, email string, code string) (bool, error) {
	key := fmt.Sprintf("otp:%s", email)
	storedCode, err := r.rdbs.Get(ctx, key).Result()

	if err == redis.Nil {
		// Key Does Not Exist
		return false, nil
	} else if err != nil {
		// Some Internal Error
		return false, fmt.Errorf("failed to retrieve OTP from redis: %w", err)
	}

	if storedCode != code {
		return false, nil
	}

	// Delete Used Key
	_ = r.rdbs.Del(ctx, key)

	return true, nil
}
