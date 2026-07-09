package authhandler

import (
	"courses/internal/auth"
	"courses/internal/config"
	"errors"
	"log"
	"net/http"
	"time"
)

type Refresh struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewRefreshHandler(l *log.Logger, authService *auth.AuthService) *Refresh {
	return &Refresh{l, authService}
}

func (m *Refresh) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_refresh_tokens")
	if err != nil || cookie.Value == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized Access"))
		return
	}

	newAccessToken, err := auth.CreateAccessToken(r.Context(), cookie.Value, m.authService)
	if err != nil {
		if errors.Is(err, auth.ErrExpiredToken) {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized Access"))
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to create access token at the moment!"))
		return
	}

	accesscookie := &http.Cookie{
		Name:     "user_access_tokens",
		Value:    newAccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(rw, accesscookie)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Access Token Cookie Added"))
}
