package models

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"html"
)

type ProfileUpdatePayload struct {
	Description string `json:"description,omitempty"`
	Password    struct {
		Old string `json:"old"`
		New string `json:"new"`
	} `json:"password,omitempty"`
}

func (payload *ProfileUpdatePayload) Sanitize() {
	payload.Description = html.EscapeString(payload.Description)
	payload.Password.Old = html.EscapeString(payload.Password.Old)
	payload.Password.New = html.EscapeString(payload.Password.New)
}

func (payload *ProfileUpdatePayload) Validate() error {
	if err := utils.CheckPassword(payload.Password.New); err != nil {
		return err
	}

	return nil
}
