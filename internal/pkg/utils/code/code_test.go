package code

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCheckCode(t *testing.T) {
	type args struct {
		code   string
		secret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				code:   "123123",
				secret: "asfdfsdvfvf",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckCode(tt.args.code, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("CheckCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
