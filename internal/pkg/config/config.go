package config

import "time"

type PayloadKey string
type RequestIdKey string

var (
	AvatarFileTypes map[string]string
	AttachFileTypes map[string]string
)

func init() {
	AvatarFileTypes = map[string]string{
		"image/jpeg": ".jpeg",
		"image/png":  ".png",
	}

	AttachFileTypes = map[string]string{
		"audio/mp4":     ".mp4",
		"audio/mpeg":    ".mp3",
		"audio/vnd.wav": ".wav",
		"image/gif":     ".gif",
		"image/jpeg":    ".jpeg",
		"image/png":     ".png",
		"video/mp4":     ".mp4",
	}
}

const (
	QrIssuer              = "YouNote"
	DefaultImagePath      = "default.jpg"
	AvatarMaxFormDataSize = 5 * 1024 * 1024
	AttachMaxFormDataSize = 30 * 1024 * 1024

	JWTLifeTime  = 24 * time.Hour
	CSRFLifeTime = 24 * time.Hour
	JwtCookie    = "YouNoteJWT"
	CsrfCookie   = "YouNoteCSRF"
	CSP          = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"

	MinUsernameLength    = 4
	MaxUsernameLength    = 12
	MinPasswordLength    = 8
	MaxPasswordLength    = 20
	PasswordAllowedExtra = "#$%&"
	SecretLength         = 6

	RedisExpirationTime = time.Minute
	MaxWrongRequests    = 5

	PayloadContextKey   PayloadKey   = "payload"
	RequestIdContextKey RequestIdKey = "request_id"

	ElasticIndexName            = "notes"
	ElasticSearchValueMinLength = 3
)
