package authhandler

import (
	"courses/internal/config"
	"courses/internal/services/auth"
	"courses/internal/services/user"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Login struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewLoginHandler(l *log.Logger, authService *auth.AuthService) *Login {
	return &Login{l, authService}
}

func (l *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	refreshToken, accessToken, userData, err := l.authService.LoginWithEmailPassword(r.Context(), creds.Email, creds.Password)

	if err != nil {
		if errors.Is(err, user.ErrUserDoesNotExist) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("User Does Not Exists!"))
			return
		} else if errors.Is(err, user.ErrWrongPassword) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Wrong Password!"))
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		errV := fmt.Sprintf("Internal Server Error! Could Not Login %v", err)
		rw.Write([]byte(errV))
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "user_refresh_tokens",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:     "user_access_tokens",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(rw).Encode(userData)
}
