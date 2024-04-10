package code

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestGenerateSecret(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test_GenerateSecret_Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteSecret := GenerateSecret()
			secret := []rune(string(byteSecret))

			assert.Equal(t, len(secret), secretLength)

			for _, sym := range secret {
				if !slices.Contains(alphabet, sym) {
					assert.Fail(t, "incorrect symbol in secret")
				}
			}
		})
	}
}
