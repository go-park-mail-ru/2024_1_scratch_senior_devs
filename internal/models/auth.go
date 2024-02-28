package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/satori/uuid"
)

const (
	minUsernameLength = 4
	maxUsernameLength = 12
	minPasswordLength = 9
	maxPasswordLength = 19
)

var (
	usernameRegexp         = regexp.MustCompile("^[0-9A-Za-z_-]+$")
	passwordRegexp         = regexp.MustCompile("^[0-9A-Za-z#$%&_-]+$")
	notLessOneLetterRegexp = regexp.MustCompile("[A-Za-z]")
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
	if !usernameRegexp.MatchString(form.Username) {
		return errors.New("username can only include symbols: A-Z, a-z, 0-9, _, - ")
	}

	if len(runedPassword) < minPasswordLength || len(runedPassword) > maxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", minPasswordLength, maxPasswordLength)
	}
	if !passwordRegexp.MatchString(form.Password) {
		return errors.New("password can only include symbols: A-Z, a-z, 0-9, #, $, %, &, _, - ")
	}
	if !notLessOneLetterRegexp.MatchString(form.Password) {
		return errors.New("password must include at least 1 letter (A-Z, a-z)")
	}

	return nil
}

type JwtPayload struct {
	Id       uuid.UUID
	Username string
}
