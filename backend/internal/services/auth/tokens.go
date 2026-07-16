package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"courses/internal/config"
	"courses/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = errors.New("Expired Token")
)

type contextKey string

const ClaimsKey contextKey = "claims"

func (a *AuthService) CreateRefreshToken(ctx context.Context, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		"type":  "Refresh",
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	err = a.RefreshRepo.SaveRefreshToken(ctx, tokenString, email)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) CreateAccessToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"].(string) != "Refresh" {
			return "", fmt.Errorf("Invalid token type")
		}

		exist, err := a.RefreshRepo.VerifyIfRefreshTokenIsLive(ctx, refreshToken, claims["email"].(string))
		if err != nil {
			return "", ErrExpiredToken
		}

		if exist {
			accessToken, err := jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				jwt.MapClaims{
					"email": claims["email"].(string),
					"exp":   time.Now().Add(time.Hour).Unix(),
					"type":  "Access",
				}).SignedString([]byte(config.JWTSecret))

			if err != nil {
				return "", errors.New("Unable to Sign Token!")
			}

			return accessToken, nil
		}
	}

	return "", ErrExpiredToken
}

func (a *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, *models.User, error) {
	accessToken, err := a.CreateAccessToken(ctx, refreshToken)
	if err != nil {
		return "", nil, err
	}

	token, _ := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	userData, err := a.User.GetUserData(ctx, email)
	if err != nil {
		return "", nil, err
	}

	return accessToken, userData, nil
}

func (a *AuthService) VerifyRefreshToken(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"].(string) == "Refresh" {
			exist, err := a.RefreshRepo.VerifyIfRefreshTokenIsLive(ctx, tokenString, claims["email"].(string))
			if err != nil {
				return false, err
			}

			return exist, nil
		}
	}

	return false, fmt.Errorf("Unauthorized")
}

func VerifyAccessToken(cntx context.Context, tokenString string) (context.Context, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return cntx, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"].(string) == "Access" {
			ctx := context.WithValue(cntx, ClaimsKey, claims)
			return ctx, nil
		}
	}

	return cntx, fmt.Errorf("Unauthorized")
}

func (a *AuthService) DeleteRefreshToken(ctx context.Context, email string) error {
	return a.RefreshRepo.DeleteRefreshToken(ctx, email)
}
