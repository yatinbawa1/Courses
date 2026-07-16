package authhandler

import (
	"courses/internal/config"
	"courses/internal/services/auth"
	"encoding/json"
	"errors"
	"fmt"
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

	newAccessToken, userData, err := m.authService.RefreshAccessToken(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, auth.ErrExpiredToken) {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized Access"))
			return
		}

		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized Access"))
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "user_access_tokens",
		Value:    newAccessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	})

	if userData == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		v := fmt.Sprintf("Unable to find user in database")
		rw.Write([]byte(v))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(userData)
}
