package authhandler

import (
	"context"
	"courses/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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
