package handlers

import (
	"courses/internal/auth"
	"courses/internal/mailer"
	"log"
	"net/http"
)

func RegisterRoutes(logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", NewHomeHandler(logger))
	mux.Handle("POST /login", NewLoginHandler(logger, authService))
	mux.Handle("POST /sign-up", NewSignUpHandler(logger, authService, mailer))

	return mux
}
