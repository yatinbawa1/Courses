package accounthandler

import (
	"courses/internal/models"
	"courses/internal/services/auth"
	"courses/internal/services/user"
	"encoding/json"
	"log"
	"net/http"
)

type UpdateUser struct {
	l           *log.Logger
	userService *user.UserService
}

func NewUpdateUserHandler(l *log.Logger, userService *user.UserService) *UpdateUser {
	return &UpdateUser{l, userService}
}

func (u *UpdateUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized"))
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unable to understand request"))
		return
	}

	user.User_id = userID

	err = u.userService.UpdateUser(r.Context(), &user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Updated"))
}
