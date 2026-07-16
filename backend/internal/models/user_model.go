package models

import (
	"github.com/google/uuid"
)

type User struct {
	User_id         uuid.UUID  `json:"user_id"`
	Username        *string    `json:"name"`
	ProfilePhotoExists bool    `json:"profile_photo_exists"`
	Email           string     `json:"email"`
}


type UserAuthCreds struct{
	Email string `json:"email"`
	Password string `json:"password"`
}