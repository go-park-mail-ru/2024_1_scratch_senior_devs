package models

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

var testConfigProfile config.UserValidationConfig

func init() {
	testLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	testConfigProfile = config.LoadConfig("../pkg/config/config.yaml", testLogger).UserValidation
}

func TestValidateProfile(t *testing.T) {
	var tests = []struct {
		name  string
		data  ProfileUpdatePayload
		isErr bool
	}{
		{
			name: "UserFormFata_ValidateSuccess",
			data: ProfileUpdatePayload{
				Description: "abc",
				Password: passwords{
					Old: "12345678a",
					New: "12345678b",
				},
			},
			isErr: false,
		},
		{
			name: "UserFormFata_ValidateFail",
			data: ProfileUpdatePayload{
				Description: "abc",
				Password: passwords{
					Old: "12345678a",
					New: "12345678",
				},
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate(testConfigAuth)

			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSanitizeProfile(t *testing.T) {
	var tests = []struct {
		name   string
		input  ProfileUpdatePayload
		output ProfileUpdatePayload
	}{
		{
			name: "Test_SanitizeProfile_Success",
			input: ProfileUpdatePayload{
				Description: "<script>alert('XSS attack')</script>",
			},
			output: ProfileUpdatePayload{
				Description: "&lt;script&gt;alert(&#39;XSS attack&#39;)&lt;/script&gt;",
			},
		},
	}

	for _, tt := range tests {
		tt.input.Sanitize()

		assert.Equal(t, tt.output, tt.input)
	}
}