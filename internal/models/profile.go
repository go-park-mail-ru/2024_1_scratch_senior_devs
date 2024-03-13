package models

import (
	"html"
)

type passwords struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type ProfileUpdatePayload struct {
	Description string    `json:"description,omitempty"`
	Password    passwords `json:"password,omitempty"`
}

func (payload *ProfileUpdatePayload) Sanitize() {
	payload.Description = html.EscapeString(payload.Description)
	payload.Password.Old = html.EscapeString(payload.Password.Old)
	payload.Password.New = html.EscapeString(payload.Password.New)
}

func (payload *ProfileUpdatePayload) Validate() error {
	if err := checkPassword(payload.Password.New); err != nil {
		return err
	}

	return nil
}
