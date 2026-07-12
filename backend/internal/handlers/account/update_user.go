package accounthandler

import (
	"courses/internal/auth"
	"courses/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

type UpdateUser struct {
	l           *log.Logger
	authService *auth.AuthService
}

func NewUpdateUserHandler(l *log.Logger, authService *auth.AuthService) *UpdateUser {
	return &UpdateUser{l, authService}
}

func (u *UpdateUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)	
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to understand request"))
		return
	}

	err = u.authService.UserRepo.SaveUser(r.Context(), &user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}	

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Updated"))
}

