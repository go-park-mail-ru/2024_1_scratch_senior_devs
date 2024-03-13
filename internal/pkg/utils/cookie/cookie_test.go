package cookie

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenTokenCookie(t *testing.T) {
	var tests = []struct {
		name    string
		token   string
		expTime time.Time
	}{
		{
			name:    "GenTokenCookie_Success",
			token:   "abc123",
			expTime: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := GenTokenCookie(tt.token, tt.expTime)
			assert.Equal(t, cookie.Name, jwt.JwtCookie)
			assert.Equal(t, cookie.Value, tt.token)
		})
	}
}

func TestDelTokenCookie(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			name: "DelTokenCookie_Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := DelTokenCookie()
			assert.Equal(t, cookie.Name, jwt.JwtCookie)
			assert.Equal(t, cookie.Value, "")
		})
	}
}
