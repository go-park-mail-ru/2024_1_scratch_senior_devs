package models

import (
	"html"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/validation"
)

type Passwords struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type ProfileUpdatePayload struct {
	Description string    `json:"description,omitempty"`
	Password    Passwords `json:"password,omitempty"`
}

func (payload *ProfileUpdatePayload) Sanitize() {
	payload.Description = html.EscapeString(payload.Description)
	payload.Password.Old = html.EscapeString(payload.Password.Old)
	payload.Password.New = html.EscapeString(payload.Password.New)
}

func (payload *ProfileUpdatePayload) Validate(cfg config.ValidationConfig) error {
	if err := validation.CheckPassword(payload.Password.New, cfg.MinPasswordLength, cfg.MaxPasswordLength, cfg.PasswordAllowedExtra); err != nil {
		return err
	}

	return nil
}
