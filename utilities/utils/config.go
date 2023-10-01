package utils

import (
	"bytes"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/Smylet/symlet-backend/resources/env"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DatabaseURI           string        `mapstructure:"DATABASE_URI"`
	DatabaseSlowThreshold time.Duration `mapstructure:"DATABASE_SLOW_THRESHOLD"`
	DatabasePoolMax       int           `mapstructure:"DATABASE_POOL_MAX"`
	DatabaseReset         bool          `mapstructure:"DATABASE_RESET"`
	DatabaseMigrate       bool          `mapstructure:"DATABASE_MIGRATE"`
	DBHost                string        `mapstructure:"DB_HOST"`
	DBPort                string        `mapstructure:"DB_PORT"`
	DBUser                string        `mapstructure:"DB_USER"`
	DBPass                string        `mapstructure:"DB_PASS"`
	DBName                string        `mapstructure:"DB_NAME"`
	SSLMode               string        `mapstructure:"SSL_MODE"`
	Environment           string        `mapstructure:"ENVIRONMENT"`
	DBSource              string        `mapstructure:"DB_SOURCE"`
	MigrationURL          string        `mapstructure:"MIGRATION_URL"`
	RedisAddress          string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress     string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress     string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName       string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress    string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword   string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	AWSRegion             string        `mapstructure:"AWS_REGION"`
	AwsAccessKeyID        string        `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey    string        `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	BasePath              string        `mapstructure:"BASE_PATH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	config_bytes := env.GetEnv()
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	err = viper.ReadConfig(bytes.NewReader(config_bytes))

	// err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
