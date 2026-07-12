package authhandler

import (
	"courses/internal/auth"
	"log"
	"net/http"
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

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Successful"))
}
