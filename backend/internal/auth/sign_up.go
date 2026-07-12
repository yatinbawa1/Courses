package auth

import (
	"context"
	"courses/internal/models"
	"errors"
	"strings"

	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	hasLower = regexp.MustCompile(`[a-z]`)
	hasUpper = regexp.MustCompile(`[A-Z]`)
	hasNum   = regexp.MustCompile(`\d`)
)

var (
	ErrUnsecurePassword       = errors.New("Password Not Secure Enough")
	ErrPasswordHashGeneration = errors.New("Unable to create a Secure Hash for Password")
	ErrWrongPassword          = errors.New("Wrong Password")
	ErrUserDoesNotExist       = errors.New("User Does Not Exist") // Used By User Repo
)

func (a *AuthService) SignUpUsingEmailAndPassword(ctx context.Context, email string, password string) error {

	email = strings.ToLower(email)

	if len(password) < 8 || !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNum.MatchString(password) {
		return ErrUnsecurePassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return ErrUnsecurePassword
	}

	var user = &models.UserAuthCreds{
		Email:    email,
		Password: string(hashedPassword),
	}

	err = a.UserRepo.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
