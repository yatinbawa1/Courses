package auth

import (
	"context"
	"courses/internal/models"
)

func (a *AuthService) LoginWithEmailPassword(ctx context.Context, email string, password string) (string, string, *models.User, error) {
	err := a.User.ValidateCredentials(ctx, email, password)
	if err != nil {
		return "", "", nil, err
	}

	userData, err := a.User.GetUserData(ctx, email)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := a.CreateRefreshToken(ctx, email, userData.User_id)
	if err != nil {
		return "", "", nil, err
	}

	accessToken, err := a.CreateAccessToken(ctx, refreshToken)
	if err != nil {
		return "", "", nil, err
	}

	return refreshToken, accessToken, userData, nil
}
