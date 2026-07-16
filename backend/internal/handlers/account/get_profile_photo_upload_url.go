package accounthandler

import (
	"courses/internal/services/user"
	"encoding/json"
	"net/http"
)

type GetProfilePhotoUploadUrl struct {
	userService *user.UserService
}

func NewGetProfilePhotoUploadUrl(userService *user.UserService) *GetProfilePhotoUploadUrl {
	return &GetProfilePhotoUploadUrl{userService}
}

func (g *GetProfilePhotoUploadUrl) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")

	resp, err := g.userService.CreatePresignedUploadURLForProfilePhoto(r.Context(), user_id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unable to create a presigned url"))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}
