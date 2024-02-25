package models

import (
	"github.com/satori/uuid"
)

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}
