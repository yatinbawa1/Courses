package handlers

import (
	"courses/internal/auth"
	"courses/internal/mailer"
	"courses/internal/middleware"
	"log"
	"net/http"
)

func RegisterRoutes(logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", middleware.CheckAuth(NewHomeHandler(logger)))
	mux.Handle("POST /login", NewLoginHandler(logger, authService))
	mux.Handle("POST /send-otp", NewSendOTPHandler(logger, authService, mailer))
	mux.Handle("POST /send-otp/verify", NewVerifyOTP(logger, authService))

	return mux
}
