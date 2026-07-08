package handlers

import (
	"context"
	"courses/internal/auth"
	"courses/internal/mailer"
	"courses/internal/models"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

// --------------------
// Login Logic
// --------------------
type Login struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewLoginHandler(l *log.Logger, authService *auth.AuthService) *Login {
	return &Login{l, authService}
}

func (l *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	token, err := l.authService.LoginWithEmailPassword(r.Context(), user.Email, user.HashedPassword)
	if err != nil {
		if errors.Is(err, auth.ErrUserDoesNotExist) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("User Does Not Exists!"))
			return
		} else if errors.Is(err, auth.ErrWrongPassword) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Wrong Password!"))
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		errV := fmt.Sprintf("Internal Server Error! Could Not Login %v", err)
		rw.Write([]byte(errV))
		return
	}

	rw.WriteHeader(http.StatusOK)
	Refreshcookie := &http.Cookie{
		Name:     "user_refresh_tokens",
		Value:    token[0],
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	Accesscookie := &http.Cookie{
		Name:     "user_access_tokens",
		Value:    token[0],
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(rw, Refreshcookie)
	http.SetCookie(rw, Accesscookie)

	rw.Write([]byte("Success"))
}

// --------------------
// Send OTP Logic	"POST /send-otp"
// --------------------

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

func (s *SendOTP) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to Parse Data"))
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

// --------------------
// Verify OTP Logic  "POST /sign-up/verify"
// --------------------

type VerifyOTPData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

type VerifyOTP struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewVerifyOTP(l *log.Logger, authService *auth.AuthService) *VerifyOTP {
	return &VerifyOTP{l, authService}
}

func (v *VerifyOTP) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user VerifyOTPData
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		ans := fmt.Sprintf("Unable to parse User Information! %s", err)
		rw.Write([]byte(ans))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	val, err := v.authService.OtpRepo.VerifyOTP(ctx, user.Email, user.OTP)
	if err != nil || !val {
		rw.WriteHeader(http.StatusBadRequest)
		ans := fmt.Sprintf("Unable to verify OTP! %s", err)
		rw.Write([]byte(ans))
		return
	}

	ctx2, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = v.authService.SignUpUsingEmailAndPassword(ctx2, user.Email, user.Password)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("Successfully Created Account"))
}
