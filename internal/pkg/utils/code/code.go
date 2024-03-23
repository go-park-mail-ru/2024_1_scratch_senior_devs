package code

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"github.com/dgryski/dgoogauth"
	"math/big"
)

var alphabet = []rune("QWETYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890")

func GenerateSecret() []byte {
	size := 30
	n := big.NewInt(int64(len(alphabet)))

	secret := make([]rune, size)
	for i := range secret {
		randomIndex, _ := rand.Int(rand.Reader, n)
		secret[i] = alphabet[randomIndex.Int64()]
	}

	return []byte(string(secret))
}

func CheckCode(code string, secret string) error {
	byteSecret := []byte(secret)

	otpConfig := &dgoogauth.OTPConfig{
		Secret:      base32.StdEncoding.EncodeToString(byteSecret),
		WindowSize:  30,
		HotpCounter: 0,
	}

	success, err := otpConfig.Authenticate(code)
	if success && err == nil {
		return nil
	} else {
		return errors.New("wrong code")
	}
}
