package models

import (
	"github.com/google/uuid"
)

type User struct {
	User_id         uuid.UUID `json:"user_id"`
	Username        string    `json:"name"`
	ProfilePhotoURL string    `json:"profile_photo_url"`
	HashedPassword  string    `json:"hashed_password"`
	Revenue         float64   `json:"revenue"`
	Email           string    `json:"email"`
}
