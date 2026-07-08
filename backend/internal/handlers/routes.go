package handlers

import (
	"courses/internal/auth"
	"courses/internal/mailer"
	"log"
	"net/http"
)

func RegisterRoutes(fileServer http.Handler, logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", fileServer)
	mux.Handle("POST /api/login", NewLoginHandler(logger, authService))
	mux.Handle("POST /api/send-otp", NewSendOTPHandler(logger, authService, mailer))
	mux.Handle("POST /api/send-otp/verify", NewVerifyOTP(logger, authService))

	return mux
}
