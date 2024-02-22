package models

import "time"

type JWTPayload struct {
	Username       string    `json:"sub"`
	ExpirationTime time.Time `json:"exp"`
}
