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
	ErrWrongPassword          = errors.New("Wrong Password")
	ErrUserDoesNotExist       = errors.New("User Does Not Exist") // Used By User Repo
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

	refreshToken, err := CreateRefreshToken(email)
	if err != nil {
		return [2]string{}, err
	}

	accessToken, err := CreateAccessToken(refreshToken)
	if err != nil {
		return [2]string{}, err
	}

	return [2]string{refreshToken, accessToken}, nil
}

func (a *AuthService) SignUpUsingEmailAndPassword(ctx context.Context, email string, password string) error {

	email = strings.ToLower(email)

	if len(password) < 8 || !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNum.MatchString(password) {
		return ErrUnsecurePassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return ErrUnsecurePassword
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
		return err
	}

	return nil
}
