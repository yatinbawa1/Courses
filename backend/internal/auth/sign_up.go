package auth

import (
	"context"
	"courses/internal/models"
	"errors"
	"strings"

	"regexp"

	"github.com/google/uuid"
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
)

func (a *AuthService) SignUpUsingEmailAndPassword(ctx context.Context, email string, password string) (string, error) {

	email = strings.ToLower(email)

	if len(password) < 8 || !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNum.MatchString(password) {
		return "", ErrUnsecurePassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", ErrUnsecurePassword
	}

	userId := uuid.New()

	var user = &models.User{
		HashedPassword:  string(hashedPassword),
		User_id:         userId,
		Email:           email,
		ProfilePhotoURL: "",
		Username:        "",
		Revenue:         0,
	}

	err = a.UserRepo.Add(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := CreateToken(email)
	if err != nil {
		return "", err
	}

	return token, nil
}
