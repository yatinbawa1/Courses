package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepo struct {
	rdbs *redis.Client
}

func NewRedisRefreshTokenRepo(rdbs *redis.Client) *RefreshTokenRepo {
	return &RefreshTokenRepo{rdbs}
}

func (r *RefreshTokenRepo) SaveRefreshToken(ctx context.Context, tokenString string, email string) error {
	key := fmt.Sprintf("refresh_token:%s", email)

	err := r.rdbs.Set(ctx, key, tokenString, time.Hour*24*30).Err()
	if err != nil {
		return fmt.Errorf("failed to Refresh Token to redis: %w", err)
	}

	return nil
}

// For Refresh Token to be verified, I need to check the passed email
// See if redis has the refresh token in it
func (r *RefreshTokenRepo) VerifyIfRefreshTokenIsLive(ctx context.Context, tokenString string, email string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s", email)
	token, err := r.rdbs.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return (token == tokenString), nil
}


func (r *RefreshTokenRepo) DeleteRefreshToken(ctx context.Context, email string) (error) {
	key := fmt.Sprintf("refresh_token:%s", email)
	_ , err :=  r.rdbs.Del(ctx, key).Result()

	return err
}
