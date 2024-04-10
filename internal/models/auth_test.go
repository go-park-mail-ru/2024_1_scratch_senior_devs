package models

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAuth(t *testing.T) {
	testConfigAuth := config.ValidationConfig{
		MinUsernameLength:    4,
		MaxUsernameLength:    12,
		MinPasswordLength:    8,
		MaxPasswordLength:    20,
		PasswordAllowedExtra: "$%&#",
		SecretLength:         6,
	}

	var tests = []struct {
		name        string
		data        UserFormData
		expectedErr error
	}{
		{
			name:        "UserFormFata_ValidateSuccess_1",
			data:        UserFormData{Username: "testuser", Password: "348gv%#332"},
			expectedErr: nil,
		},
		{
			name:        "UserFormFata_ValidateSuccess_2",
			data:        UserFormData{Username: "testuser2", Password: "nv48392fh$"},
			expectedErr: nil,
		},
		{
			name:        "UserFormFata_ValidateFail_1",
			data:        UserFormData{Username: "+74951234567", Password: "nv48392fh$"},
			expectedErr: errors.New("username can only include symbols: A-Z, a-z, 0-9"),
		},
		{
			name:        "UserFormFata_ValidateFail_2",
			data:        UserFormData{Username: "7495123456789", Password: "nv48392fh$"},
			expectedErr: fmt.Errorf("username length must be from %d to %d characters", testConfigAuth.MinUsernameLength, testConfigAuth.MaxUsernameLength),
		},
		{
			name:        "UserFormFata_ValidateFail_3",
			data:        UserFormData{Username: "74951234567", Password: "fn4839vjn8309jn80c39j23hfv93n309h4v"},
			expectedErr: fmt.Errorf("password length must be from %d to %d characters", testConfigAuth.MinPasswordLength, testConfigAuth.MaxPasswordLength),
		},
		{
			name:        "UserFormFata_ValidateFail_4",
			data:        UserFormData{Username: "74951234567", Password: "cn39 3297yth2"},
			expectedErr: errors.New("password can only include symbols: A-Z, a-z, 0-9, #, $, %, &"),
		},
		{
			name:        "UserFormFata_ValidateFail_5",
			data:        UserFormData{Username: "74951234567", Password: "42368723632"},
			expectedErr: errors.New("password must include at least 1 letter (A-Z, a-z)"),
		},
		{
			name:        "UserFormFata_ValidateFail_6",
			data:        UserFormData{Username: "testuser", Password: "348gv%#332", Code: "12345"},
			expectedErr: fmt.Errorf("secret length must be %d", testConfigAuth.SecretLength),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate(testConfigAuth)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
