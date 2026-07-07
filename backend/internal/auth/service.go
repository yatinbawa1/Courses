package auth

import (
	"context"
	"courses/internal/models"
	"courses/internal/repository"
)

type AuthService struct {
	UserRepo UserDataRepo
	OtpRepo  repository.OTPRepo
}

type UserDataRepo interface {
	Add(ctx context.Context, use *models.User) error
}

func NewAuthService(userRepo UserDataRepo, otpRepo repository.OTPRepo) *AuthService {
	return &AuthService{userRepo, otpRepo}
}
