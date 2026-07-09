package auth

import (
	"context"
	"courses/internal/models"
	"time"
)

type AuthService struct {
	UserRepo    UserDataRepo
	OtpRepo     OTPRepo
	RefreshRepo RefreshTokenRepo
}

type UserDataRepo interface {
	Add(ctx context.Context, use *models.User) error
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	GetPasswordForEmail(ctx context.Context, email string) ([]byte, error)
}

type OTPRepo interface {
	SaveOTP(ctx context.Context, email string, code string, tll time.Duration) error
	VerifyOTP(ctx context.Context, email string, code string) (bool, error)
}

type RefreshTokenRepo interface {
	SaveRefreshToken(ctx context.Context, tokenString string, email string) error
	VerifyIfRefreshTokenIsLive(ctx context.Context, tokenString string, email string) (bool, error)
}

func NewAuthService(userRepo UserDataRepo, otpRepo OTPRepo, refreshTokenRepo RefreshTokenRepo) *AuthService {
	return &AuthService{userRepo, otpRepo, refreshTokenRepo}
}
