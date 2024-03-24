package models

import (
	"encoding/json"
	"time"

	"github.com/satori/uuid"
)

type Secret string

func (secret Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(secret != "")
}

type User struct {
	Id           uuid.UUID `json:"id"`
	Description  string    `json:"description,omitempty"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"create_time"`
	ImagePath    string    `json:"image_path"`
	SecondFactor Secret    `json:"second_factor"`
}

// ================================================================
// only swagger examples

type UserForSwagger struct {
	Id           uuid.UUID `json:"id"`
	Description  string    `json:"description,omitempty"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"create_time"`
	ImagePath    string    `json:"image_path"`
	SecondFactor bool      `json:"second_factor"`
}
