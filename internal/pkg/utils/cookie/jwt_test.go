package cookie

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

var testConfig config.JwtConfig

func init() {
	testLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	testConfig = config.LoadConfig("", testLogger).AuthHandler.Jwt
}

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
			jwtCookie := GenJwtTokenCookie(tt.token, tt.expTime, testConfig)
			assert.Equal(t, jwtCookie.Name, testConfig.JwtCookie)
			assert.Equal(t, jwtCookie.Value, tt.token)
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
			jwtCookie := DelJwtTokenCookie(testConfig)
			assert.Equal(t, jwtCookie.Name, testConfig.JwtCookie)
			assert.Equal(t, jwtCookie.Value, "")
		})
	}
}
