package authhandler

import (
	"courses/internal/auth"
	"log"
	"net/http"
	"time"
)

type Logout struct {
	l *log.Logger
	authService *auth.AuthService
}

func NewLogoutHandler(l *log.Logger, a *auth.AuthService) *Logout {
	return &Logout{l,a}
}

func (l *Logout) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	userEmail := r.PathValue("user_email")
	err := l.authService.RefreshRepo.DeleteRefreshToken(r.Context(),userEmail)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to delete refresh token"))
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "user_refresh_tokens",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:     "user_access_tokens",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Successful"))
}
