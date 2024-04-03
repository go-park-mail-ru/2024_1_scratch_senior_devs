package config

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type PayloadKey string
type RequestIdKey string

type Config struct {
	AuthHandler    AuthHandlerConfig    `yaml:"auth_handler"`
	AuthUsecase    AuthUsecaseConfig    `yaml:"auth_usecase"`
	Blocker        BlockerConfig        `yaml:"blocker"`
	UserValidation UserValidationConfig `yaml:"user_validation"`
	Attach         AttachConfig         `yaml:"attach"`
}

type AuthHandlerConfig struct {
	QrIssuer              string            `yaml:"qr_issuer"`
	AvatarMaxFormDataSize int               `yaml:"avatar_max_form_data_size"`
	AvatarFileTypes       map[string]string `yaml:"avatar_file_types"`
	Jwt                   JwtConfig         `yaml:"jwt"`
	Csrf                  CsrfConfig        `yaml:"csrf"`
}

type AuthUsecaseConfig struct {
	DefaultImagePath string        `yaml:"default_image_path"`
	JWTLifeTime      time.Duration `yaml:"jwt_lifetime"`
}

type JwtConfig struct {
	JwtCookie string `yaml:"jwt_cookie"`
}
type CsrfConfig struct {
	CsrfCookie   string        `yaml:"csrf_cookie"`
	CSRFLifeTime time.Duration `yaml:"csrf_lifetime"`
}

type AttachConfig struct {
	AttachMaxFormDataSize int64             `yaml:"attach_max_form_data_size"`
	AttachFileTypes       map[string]string `yaml:"attach_file_types"`
}
type BlockerConfig struct {
	RedisExpirationTime time.Duration `yaml:"redis_expiration_time"`
	MaxWrongRequests    int           `yaml:"max_wrong_requests"`
}
type UserValidationConfig struct {
	MinUsernameLength    int    `yaml:"min_username_length"`
	MaxUsernameLength    int    `yaml:"max_username_length"`
	MinPasswordLength    int    `yaml:"min_password_length"`
	MaxPasswordLength    int    `yaml:"max_password_length"`
	PasswordAllowedExtra string `yaml:"password_allowed_extra"`
	SecretLength         int    `yaml:"secret_length"`
}

const (
	CSP                              = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"
	PayloadContextKey   PayloadKey   = "payload"
	RequestIdContextKey RequestIdKey = "request_id"
)

var cfg *Config
var once sync.Once

func LoadConfig(path string, logger *slog.Logger) *Config {
	once.Do(func() {
		cfg = &Config{}
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Cant open config file")
			return
		}
		defer file.Close()
		d := yaml.NewDecoder(file)
		if err := d.Decode(cfg); err != nil {
			fmt.Println("Cant decode config", err.Error())
			return
		}

	})
	return cfg
}
