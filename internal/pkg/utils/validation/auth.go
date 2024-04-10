package validation

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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

func checkPasswordAllowed(value []rune, PasswordAllowedExtra string) bool {
	for _, sym := range value {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) && !strings.Contains(PasswordAllowedExtra, string(sym)) {
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

func CheckUsername(username string, MinUsernameLength int, MaxUsernameLength int) error {
	runedUsername := []rune(username)

	if len(runedUsername) < MinUsernameLength || len(runedUsername) > MaxUsernameLength {
		return fmt.Errorf("username length must be from %d to %d characters", MinUsernameLength, MaxUsernameLength)
	}

	if !checkUsernameAllowed(runedUsername) {
		return errors.New("username can only include symbols: A-Z, a-z, 0-9")
	}

	return nil
}

func CheckPassword(password string, MinPasswordLength int, MaxPasswordLength int, PasswordAllowedWxtra string) error {
	runedPassword := []rune(password)

	if len(runedPassword) < MinPasswordLength || len(runedPassword) > MaxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", MinPasswordLength, MaxPasswordLength)
	}

	if !checkPasswordAllowed(runedPassword, PasswordAllowedWxtra) {
		return errors.New("password can only include symbols: A-Z, a-z, 0-9, #, $, %, &")
	}

	if !checkPasswordRequired(runedPassword) {
		return errors.New("password must include at least 1 letter (A-Z, a-z)")
	}

	return nil
}

func CheckSecret(secret string, SecretLength int) error {
	runedSecret := []rune(secret)

	if len(runedSecret) == 0 {
		return nil
	}

	if len(runedSecret) != SecretLength {
		return fmt.Errorf("secret length must be %d", SecretLength)
	}

	for _, sym := range runedSecret {
		if !unicode.IsDigit(sym) {
			return fmt.Errorf("secret must contain only digits")
		}
	}

	return nil
}
