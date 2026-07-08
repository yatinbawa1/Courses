package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"courses/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = errors.New("Expired Token")
)

func CreateRefreshToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 15).Unix(),
		"type":  "Refresh",
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateAccessToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
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
	} else {
		return "", ErrExpiredToken
	}
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
			ctx := context.WithValue(cntx, "claims", claims)
			return ctx, nil
		}
	}

	return cntx, fmt.Errorf("Unauthorized")
}
