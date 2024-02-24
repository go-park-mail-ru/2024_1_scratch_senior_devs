package models

import (
	"time"

	"github.com/satori/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Description  string    `json:"description"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"create_time"`
	ImagePath    string    `json:"image_path"`
}
