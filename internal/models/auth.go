package models

import (
	"fmt"
	"github.com/satori/uuid"
)

const (
	minUsernameLength = 2
	maxUsernameLength = 12
	minPasswordLength = 8
	maxPasswordLength = 64
)

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (form *UserFormData) Validate() error {
	runedUsername := []rune(form.Username)
	runedPassword := []rune(form.Password)

	if len(runedUsername) < minUsernameLength || len(runedUsername) > maxUsernameLength {
		return fmt.Errorf("username length must be from %d to %d characters", minUsernameLength, maxUsernameLength)
	}

	if len(runedPassword) < minPasswordLength || len(runedUsername) > maxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", minPasswordLength, maxPasswordLength)
	}

	return nil
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}
