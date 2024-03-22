package qrcode

import (
	"math/rand"
	"unicode/utf8"
)

const alphabet = "QWETYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

func GenerateSecret() []byte {
	n := utf8.RuneCountInString(alphabet)

	var secret []byte
	for i := 0; i < 30; i++ {
		secret = append(secret, alphabet[rand.Int()%n])
	}

	return secret
}
