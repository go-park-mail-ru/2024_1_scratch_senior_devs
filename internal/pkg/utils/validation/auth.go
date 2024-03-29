package validation

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
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

func checkPasswordAllowed(value []rune) bool {
	for _, sym := range value {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) && !strings.Contains(config.PasswordAllowedExtra, string(sym)) {
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

	if len(runedUsername) < config.MinUsernameLength || len(runedUsername) > config.MaxUsernameLength {
		return fmt.Errorf("username length must be from %d to %d characters", config.MinUsernameLength, config.MaxUsernameLength)
	}

	if !checkUsernameAllowed(runedUsername) {
		return errors.New("username can only include symbols: A-Z, a-z, 0-9")
	}

	return nil
}

func CheckPassword(password string) error {
	runedPassword := []rune(password)

	if len(runedPassword) < config.MinPasswordLength || len(runedPassword) > config.MaxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", config.MinPasswordLength, config.MaxPasswordLength)
	}

	if !checkPasswordAllowed(runedPassword) {
		return errors.New("password can only include symbols: A-Z, a-z, 0-9, #, $, %, &")
	}

	if !checkPasswordRequired(runedPassword) {
		return errors.New("password must include at least 1 letter (A-Z, a-z)")
	}

	return nil
}

func CheckSecret(secret string) error {
	runedSecret := []rune(secret)

	if len(runedSecret) == 0 {
		return nil
	}

	if len(runedSecret) != config.SecretLength {
		return fmt.Errorf("secret length must be %d", config.SecretLength)
	}

	for _, sym := range runedSecret {
		if !unicode.IsDigit(sym) {
			return fmt.Errorf("secret must contain only digits")
		}
	}

	return nil
}
