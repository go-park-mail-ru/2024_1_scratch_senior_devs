package models

import (
	"time"
)

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreateTime   time.Time `json:"create_time"`
	ImagePath    string    `json:"image_path"`
}
