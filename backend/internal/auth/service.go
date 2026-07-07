package auth

import (
	"context"
	"courses/internal/models"
)

type AuthService struct {
	userRepo UserDataRepo
}

type UserDataRepo interface {
	Add(ctx context.Context, use *models.User) error
}

func NewAuthService(userRepo UserDataRepo) *AuthService {
	return &AuthService{userRepo}
}
