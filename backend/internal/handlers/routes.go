package handlers

import (
	"courses/internal/auth"
	"log"
	"net/http"

	"github.com/resend/resend-go/v3"
)

func RegisterRoutes(logger *log.Logger, authService *auth.AuthService, mailClient *resend.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", NewHomeHandler(logger))
	mux.Handle("POST /login", NewLoginHandler(logger, authService))
	mux.Handle("POST /sign-up", NewSignUpHandler(logger, authService, mailClient))

	return mux
}
