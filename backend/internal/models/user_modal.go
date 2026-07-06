package models

import (
	"github.com/google/uuid"
)

type User struct {
	User_id         uuid.UUID `json:"user_id"`
	Username        string    `json:"username"`
	ProfilePhotoURL string    `json:"profile_photo_url"`
	HashedPassword  string    `json:"hashed_password"`
	Revenue         int       `json:"revenue"`
}
