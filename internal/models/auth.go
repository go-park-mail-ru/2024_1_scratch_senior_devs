package models

import (
	"time"
)

type JWTPayload struct {
	username string
	lifeTime time.Duration
}

func (payload *JWTPayload) GenJWT() {
	// expTime := time.Now().Add(payload.lifeTime)
}

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"-"`
}
