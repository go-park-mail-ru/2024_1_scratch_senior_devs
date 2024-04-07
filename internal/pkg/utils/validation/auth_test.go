package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth(t *testing.T) {
	secretLength := 6

	tests := []struct {
		name   string
		secret string
		isErr  bool
	}{
		{
			name:   "TestAuth_Success_1",
			secret: "123456",
			isErr:  false,
		},
		{
			name:   "TestAuth_Success_2",
			secret: "",
			isErr:  false,
		},
		{
			name:   "TestAuth_Fail_1",
			secret: "12345",
			isErr:  true,
		},
		{
			name:   "TestAuth_Fail_2",
			secret: "12345a",
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSecret(tt.secret, secretLength)

			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
