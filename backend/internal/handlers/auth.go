package handlers

import (
	"courses/internal/auth"
	"courses/internal/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/resend/resend-go/v3"
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
	mailClient  *resend.Client
}

func NewSignUpHandler(l *log.Logger, authService *auth.AuthService, mailClient *resend.Client) *SignUp {
	return &SignUp{l, authService, mailClient}
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
