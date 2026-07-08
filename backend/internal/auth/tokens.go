package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"courses/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(cntx context.Context, tokenString string) (context.Context, error) {
	log.Printf("Verifying Token: %s", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})
	log.Printf("Verifying Token: Token Parsed errors: %w", err)

	if err != nil {
		return cntx, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx := context.WithValue(cntx, "claims", claims)
		return ctx, nil
	}

	return cntx, fmt.Errorf("Unauthorized")
}
