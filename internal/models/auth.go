package models

import (
	"github.com/satori/uuid"
)

type PayloadKey string

const PayloadContextKey PayloadKey = "payload"

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (form *UserFormData) Validate() error {
	if err := checkUsername(form.Username); err != nil {
		return err
	}

	if err := checkPassword(form.Password); err != nil {
		return err
	}

	return nil
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}

type QrPayload struct {
	Username string `json:"username"`
}

type OtpPayload struct {
	Username string `json:"username"`
	QrCode   string `json:"qr_code"`
}
