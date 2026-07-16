package auth

import (
	"context"
)

func (a *AuthService) SignUpUsingEmailAndPassword(ctx context.Context, email string, password string) error {
	return a.User.SignUpUser(ctx, email, password)
}
