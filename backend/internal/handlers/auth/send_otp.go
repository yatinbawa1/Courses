package authhandler

import (
	"context"
	"courses/internal/auth"
	"courses/internal/mailer"
	"courses/internal/models"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
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

type SendOTP struct {
	l           *log.Logger
	authService *auth.AuthService
	mailClient  mailer.MailSender
}

func NewSendOTPHandler(l *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *SendOTP {
	return &SendOTP{l, authService, mailer}
}
func VerifyIfEmailPasswordAreOK(email string, password string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if len(password) <= 8 {
		return false
	}

	var hasUpper, hasLower, hasSymbol bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}

		if hasUpper && hasLower && hasSymbol {
			return true
		}
	}

	return hasUpper && hasLower && hasSymbol
}
func (s *SendOTP) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to Parse Data"))
		return
	}

	if !VerifyIfEmailPasswordAreOK(user.Email, user.HashedPassword) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email OR Password Not In Format"))
		return
	}
	exist, err := s.authService.UserRepo.CheckIfEmailExists(r.Context(), user.Email)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable To Check User Email In Database"))
		return
	}

	if exist {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Email Already Exists! Kindly Use Another Email"))
		return
	}

	otp := generateOTP()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = s.authService.OtpRepo.SaveOTP(ctx, user.Email, otp, time.Second*3600)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		ans := fmt.Sprintf("Unable to Save OTP in Redis! %s", err)
		rw.Write([]byte(ans))
		return
	}

	go func(email, otp string) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err = s.mailClient.SendOTPMail(ctx, &user, otp)
		if err != nil {
			s.l.Printf("Failed To Send Email To %s: %v", email, err)
		}
	}(user.Email, otp)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OTP Sent!"))
}
