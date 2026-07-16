package auth

import (
	"context"
	"courses/internal/models"
	"time"
)

type UserDataProvider interface {
	ValidateCredentials(ctx context.Context, email string, password string) error
	SignUpUser(ctx context.Context, email string, password string) error
	GetUserData(ctx context.Context, email string) (*models.User, error)
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
}

type AuthService struct {
	OtpRepo     OTPRepo
	RefreshRepo RefreshTokenRepo
	User        UserDataProvider
}

type OTPRepo interface {
	SaveOTP(ctx context.Context, email string, code string, tll time.Duration) error
	GetOTPForUser(ctx context.Context, email string) (string, error)
	DeleteOTPForUser(ctx context.Context, email string) error
}

type RefreshTokenRepo interface {
	SaveRefreshToken(ctx context.Context, tokenString string, email string) error
	VerifyIfRefreshTokenIsLive(ctx context.Context, tokenString string, email string) (bool, error)
	DeleteRefreshToken(ctx context.Context, email string) error
}

func NewAuthService(otpRepo OTPRepo, refreshTokenRepo RefreshTokenRepo, user UserDataProvider) *AuthService {
	return &AuthService{otpRepo, refreshTokenRepo, user}
}
