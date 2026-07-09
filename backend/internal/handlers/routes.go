package handlers

import (
	"courses/internal/auth"
	authhandler "courses/internal/handlers/auth"
	"courses/internal/mailer"
	"log"
	"net/http"
)

func RegisterRoutes(fileServer http.Handler, logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", fileServer)
	mux.Handle("GET /api/auth/refresh", authhandler.NewRefreshHandler(logger, authService))
	mux.Handle("POST /api/auth/login", authhandler.NewLoginHandler(logger, authService))
	mux.Handle("POST /api/auth/send-otp", authhandler.NewSendOTPHandler(logger, authService, mailer))
	mux.Handle("POST /api/auth/send-otp/verify", authhandler.NewVerifyOTP(logger, authService))

	return mux
}
