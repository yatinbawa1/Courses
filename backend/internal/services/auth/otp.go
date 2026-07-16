package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"time"
	"unicode"
)

func generateOTP() string {
	max := big.NewInt(900000)
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return ""
	}
	otp := randomNum.Int64() + 100000
	return fmt.Sprintf("%d", otp)
}

func validateEmailAndPasswordFormat(email string, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("Invalid email format")
	}

	if len(password) < 8 {
		return fmt.Errorf("Password too short")
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("Password must contain upper, lower, and digit")
	}

	return nil
}

func (a *AuthService) SendOTP(ctx context.Context, email string, password string) (string, error) {
	if err := validateEmailAndPasswordFormat(email, password); err != nil {
		return "", err
	}

	exists, err := a.User.CheckIfEmailExists(ctx, email)
	if err != nil {
		return "", fmt.Errorf("Unable to check email in database")
	}

	if exists {
		return "", fmt.Errorf("Email already exists")
	}

	otp := generateOTP()

	err = a.OtpRepo.SaveOTP(ctx, email, otp, time.Second*3600)
	if err != nil {
		return "", fmt.Errorf("Unable to save OTP")
	}

	return otp, nil
}

func (a *AuthService) VerifyOTPAndSignUp(ctx context.Context, email string, password string, otp string) error {
	val, err := a.OtpRepo.GetOTPForUser(ctx, email)
	if err != nil || val != otp {
		return fmt.Errorf("Unable to verify OTP")
	}

	a.OtpRepo.DeleteOTPForUser(ctx, email)
	return a.SignUpUsingEmailAndPassword(ctx, email, password)
}
