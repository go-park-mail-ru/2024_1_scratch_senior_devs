package models

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
)

type PayloadKey string

const PayloadContextKey PayloadKey = "payload"

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (form *UserFormData) Validate() error {
	if err := utils.CheckUsername(form.Username); err != nil {
		return err
	}

	if err := utils.CheckPassword(form.Password); err != nil {
		return err
	}

	return nil
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}
