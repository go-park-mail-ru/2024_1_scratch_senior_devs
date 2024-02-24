package models

import (
	"time"
	
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"create_time"`
	ImagePath    string    `json:"image_path"`
}
