package models

import "time"

type JWTPayload struct {
	Username       string        `json:"sub"`
	ExpirationTime time.Duration `json:"exp"`
}

type SignUpForm struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
