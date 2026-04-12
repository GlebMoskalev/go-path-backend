package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServiceConfig     `mapstructure:",squash"`
	JWT      JWTConfig         `mapstructure:",squash"`
	Google   GoogleOAuthConfig `mapstructure:",squash"`
	Database DatabaseConfig    `mapstructure:",squash"`
	Redis    RedisConfig       `mapstructure:",squash"`
	Sandbox  SandboxConfig     `mapstructure:",squash"`
	AIConfig AIConfig          `mapstructure:",squash"`
	Env      string            `mapstructure:"ENV"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type GoogleOAuthConfig struct {
	ClientID           string `mapstructure:"GOOGLE_CLIENT_ID"`
	ClientSecret       string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	CallbackURL        string `mapstructure:"GOOGLE_CALLBACK_URL"`
	FrontedCallbackURL string `mapstructure:"FRONTEND_CALLBACK_URL"`
	UserInfoURL        string `mapstructure:"GOOGLE_USER_INFO_URL"`
}

type ServiceConfig struct {
	Host               string        `mapstructure:"SERVER_HOST"`
	Port               string        `mapstructure:"SERVER_PORT"`
	ShutdownTimeoutStr string        `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
	ShutdownTimeout    time.Duration `mapstructure:"-"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"JWT_SECRET"`
	AccessTTLStr  string        `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTTLStr string        `mapstructure:"JWT_REFRESH_TTL"`
	AccessTTL     time.Duration `mapstructure:"-"`
	RefreshTTL    time.Duration `mapstructure:"-"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type SandboxConfig struct {
	Image      string        `mapstructure:"SANDBOX_IMAGE"`
	Timeout    time.Duration `mapstructure:"-"`
	TimeoutStr string        `mapstructure:"SANDBOX_TIMEOUT"`
	Memory     int64         `mapstructure:"SANDBOX_MEMORY"`
	NanoCPUs   int64         `mapstructure:"SANDBOX_NANO_CPUS"`
}

type AIConfig struct {
	ApiKey              string  `mapstructure:"AI_API_KEY"`
	ApiUrl              string  `mapstructure:"AI_API_URL"`
	ModelPassedTests    string  `mapstructure:"AI_MODEL_PASSED_TESTS"`
	MaxTokensTask       int     `mapstructure:"AI_MAX_TOKENS_TASK"`
	MaxTokensProject    int     `mapstructure:"AI_MAX_TOKENS_PROJECT"`
	Temperature         float32 `mapstructure:"AI_TEMPERATURE"`
	TopP                float32 `mapstructure:"AI_TOP_P"`
	SystemPromptTask    string  `mapstructure:"AI_SYSTEM_PROMPT_TASK"`
	UserPromptTask      string  `mapstructure:"AI_USER_PROMPT_TASK"`
	SystemPromptProject string  `mapstructure:"AI_SYSTEM_PROMPT_PROJECT"`
	UserPromptProject   string  `mapstructure:"AI_USER_PROMPT_PROJECT"`
	SystemPromptError   string  `mapstructure:"AI_SYSTEM_PROMPT_ERROR"`
	UserPromptError     string  `mapstructure:"AI_USER_PROMPT_ERROR"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	accessTTL, err := time.ParseDuration(cfg.JWT.AccessTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}
	cfg.JWT.AccessTTL = accessTTL

	refreshTTL, err := time.ParseDuration(cfg.JWT.RefreshTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}
	cfg.JWT.RefreshTTL = refreshTTL

	shutdownTimeout, err := time.ParseDuration(cfg.Server.ShutdownTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.Server.ShutdownTimeout = shutdownTimeout

	sandboxTimeout, err := time.ParseDuration(cfg.Sandbox.TimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SANDBOX_TIMEOUT: %w", err)
	}
	cfg.Sandbox.Timeout = sandboxTimeout

	return &cfg, nil
}
