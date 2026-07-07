package handlers

import (
	"courses/internal/auth"
	"courses/internal/mailer"
	"courses/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

// --------------------
// Login Logic
// --------------------
type Login struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewLoginHandler(l *log.Logger, authService *auth.AuthService) *Login {
	return &Login{l, authService}
}

func (l *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}

// --------------------
// Sign Up Logic
// --------------------

type SignUp struct {
	l           *log.Logger
	authService *auth.AuthService
	mailClient  mailer.MailSender
}

func NewSignUpHandler(l *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *SignUp {
	return &SignUp{l, authService, mailer}
}

func (s *SignUp) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to parse User Information!"))
		return
	}

	token, err := s.authService.SignUpUsingEmailAndPassword(r.Context(), user.Email, user.HashedPassword)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(token))
}
