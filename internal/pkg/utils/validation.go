package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	minUsernameLength    = 4
	maxUsernameLength    = 12
	minPasswordLength    = 8
	maxPasswordLength    = 20
	passwordAllowedExtra = "#$%&"
)

func isEnglishLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func checkUsernameAllowed(value []rune) bool {
	for _, sym := range value {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) {
			return false
		}
	}
	return true
}

func checkPasswordAllowed(value []rune) bool {
	for _, sym := range value {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) && !strings.Contains(passwordAllowedExtra, string(sym)) {
			return false
		}
	}
	return true
}

func checkPasswordRequired(value []rune) bool {
	for _, sym := range value {
		if isEnglishLetter(sym) {
			return true
		}
	}
	return false
}

func CheckUsername(username string) error {
	runedUsername := []rune(username)

	if len(runedUsername) < minUsernameLength || len(runedUsername) > maxUsernameLength {
		return fmt.Errorf("username length must be from %d to %d characters", minUsernameLength, maxUsernameLength)
	}

	if !checkUsernameAllowed(runedUsername) {
		return errors.New("username can only include symbols: A-Z, a-z, 0-9")
	}

	return nil
}

func CheckPassword(password string) error {
	runedPassword := []rune(password)

	if len(runedPassword) < minPasswordLength || len(runedPassword) > maxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", minPasswordLength, maxPasswordLength)
	}

	if !checkPasswordAllowed(runedPassword) {
		return errors.New("password can only include symbols: A-Z, a-z, 0-9, #, $, %, &")
	}

	if !checkPasswordRequired(runedPassword) {
		return errors.New("password must include at least 1 letter (A-Z, a-z)")
	}

	return nil
}
