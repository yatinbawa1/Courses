package authhandler

import (
	"context"
	"courses/internal/models"
	"courses/internal/services/auth"
	"courses/internal/services/mailer"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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

	var user models.UserAuthCreds
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to Parse Data"))
		return
	}

	otp, err := s.authService.SendOTP(r.Context(), user.Email, user.Password)
	if err != nil {
		msg := err.Error()
		switch {
		case msg == "Email already exists":
			rw.WriteHeader(http.StatusConflict)
			rw.Write([]byte("Email Already Exists! Kindly Use Another Email"))
		case msg == "Invalid email format" || msg == "Password too short" || msg == "Password must contain upper, lower, and digit":
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Email OR Password Not In Format"))
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Unable To Process Request"))
		}
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
