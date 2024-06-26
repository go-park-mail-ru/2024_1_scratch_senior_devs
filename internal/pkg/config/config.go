package config

import (
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type PayloadKey string
type RequestIdKey string
type LoggerKey string

type Config struct {
	Main        MainConfig        `yaml:"main"`
	AuthHandler AuthHandlerConfig `yaml:"auth_handler"`
	AuthUsecase AuthUsecaseConfig `yaml:"auth_usecase"`
	Blocker     BlockerConfig     `yaml:"blocker"`
	Validation  ValidationConfig  `yaml:"validation"`
	Attach      AttachConfig      `yaml:"attach"`
	Elastic     ElasticConfig     `yaml:"elastic"`
	Grpc        GrpcConfig        `yaml:"grpc"`
	Hub         HubConfig         `yaml:"hub"`
	Constraints ConstraintsConfig `yaml:"constraints"`
}

type MainConfig struct {
	Port              string        `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
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

type ValidationConfig struct {
	MinUsernameLength    int    `yaml:"min_username_length"`
	MaxUsernameLength    int    `yaml:"max_username_length"`
	MinPasswordLength    int    `yaml:"min_password_length"`
	MaxPasswordLength    int    `yaml:"max_password_length"`
	PasswordAllowedExtra string `yaml:"password_allowed_extra"`
	SecretLength         int    `yaml:"secret_length"`
}

type ElasticConfig struct {
	ElasticIndexName            string `yaml:"elastic_index_name"`
	ElasticSearchValueMinLength int    `yaml:"elastic_search_value_min_length"`
}

type GrpcConfig struct {
	AuthPort        string `yaml:"auth_port"`
	AuthIP          string `yaml:"auth_ip"`
	AuthMetricsPort string `yaml:"auth_metrics_port"`
	NotePort        string `yaml:"note_port"`
	NoteIP          string `yaml:"note_ip"`
	NoteMetricsPort string `yaml:"note_metrics_port"`
}

type HubConfig struct {
	Period   time.Duration `yaml:"period"`
	CacheTtl time.Duration `yaml:"cache_ttl"`
}

type ConstraintsConfig struct {
	MaxSubnotes      int `yaml:"max_subnotes"`
	MaxDepth         int `yaml:"max_depth"`
	MaxCollaborators int `yaml:"max_collaborators"`
	MaxTags          int `yaml:"max_tags"`
}

const (
	PayloadContextKey   PayloadKey   = "payload"
	RequestIdContextKey RequestIdKey = "request_id"
	LoggerContextKey    LoggerKey    = "logger"
)

func LoadConfig(path string, logger *slog.Logger) *Config {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		logger.Error("Cant open config file: " + err.Error())
		return cfg
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(cfg); err != nil {
		logger.Error("Cant decode config: " + err.Error())
		return cfg
	}

	logger.Info("Successfully loaded config")
	return cfg
}
