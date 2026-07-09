package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) LoginWithEmailPassword(ctx context.Context, email string, password string) ([2]string, error) {
	pass, err := a.UserRepo.GetPasswordForEmail(ctx, email)
	if err != nil {
		return [2]string{}, err
	}

	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return [2]string{}, ErrWrongPassword
	}

	refreshToken, err := CreateRefreshToken(ctx, email, a)
	if err != nil {
		return [2]string{}, err
	}

	accessToken, err := CreateAccessToken(ctx, refreshToken, a)
	if err != nil {
		return [2]string{}, err
	}

	return [2]string{refreshToken, accessToken}, nil
}
