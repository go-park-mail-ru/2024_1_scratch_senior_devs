package models

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/validation"
	"github.com/satori/uuid"
)

type PayloadKey string

const PayloadContextKey PayloadKey = "payload"

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"`
}

func (form *UserFormData) Validate() error {
	if err := validation.CheckUsername(form.Username); err != nil {
		return err
	}

	if err := validation.CheckPassword(form.Password); err != nil {
		return err
	}

	if err := validation.CheckSecret(form.Code); err != nil {
		return err
	}

	return nil
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}

// ================================================================
// only swagger examples

type SignUpPayloadForSwagger struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
