package authhandler

import (
	"courses/internal/auth"
	"courses/internal/config"
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
	
	newAccessToken, err, email := auth.CreateAccessToken(r.Context(), cookie.Value, m.authService)
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
	userData, err := m.authService.UserRepo.GetUserData(r.Context(), email)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		v := fmt.Sprintf("Unable to find user in database %s", err)
		rw.Write([]byte(v))
		return
	}

	rw.Header().Set("Content-Type","application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(userData)
}
