package models

import "time"

type JWTPayload struct {
	Username       string        `json:"sub"`
	ExpirationTime time.Duration `json:"exp"`
}
