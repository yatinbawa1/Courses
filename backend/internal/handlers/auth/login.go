package authhandler

import (
	"courses/internal/auth"
	"courses/internal/config"
	"courses/internal/models"
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

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	token, err := l.authService.LoginWithEmailPassword(r.Context(), user.Email, user.HashedPassword)

	if err != nil {
		if errors.Is(err, auth.ErrUserDoesNotExist) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("User Does Not Exists!"))
			return
		} else if errors.Is(err, auth.ErrWrongPassword) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Wrong Password!"))
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		errV := fmt.Sprintf("Internal Server Error! Could Not Login %v", err)
		rw.Write([]byte(errV))
		return
	}

	Refreshcookie := &http.Cookie{
		Name:     "user_refresh_tokens",
		Value:    token[0],
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	}

	Accesscookie := &http.Cookie{
		Name:     "user_access_tokens",
		Value:    token[1],
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIES,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(rw, Refreshcookie)
	http.SetCookie(rw, Accesscookie)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Success"))
}
