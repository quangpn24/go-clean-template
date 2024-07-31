package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Empty = new(Config)

type Config struct {
	AppEnv            string `envconfig:"APP_ENV"`
	Port              int    `envconfig:"PORT"`
	SentryDSN         string `envconfig:"SENTRY_DSN"`
	AllowOrigins      string `envconfig:"ALLOW_ORIGINS"`
	CognitoIssuer     string `envconfig:"COGNITO_ISSUER"`
	CognitoURLGetJWKS string `envconfig:"COGNITO_URL_GET_JWKS"`
	UserPoolID        string `envconfig:"USER_POOL_ID"`

	DB struct {
		Name      string `envconfig:"DB_NAME"`
		Host      string `envconfig:"DB_HOST"`
		Port      int    `envconfig:"DB_PORT"`
		User      string `envconfig:"DB_USER"`
		Pass      string `envconfig:"DB_PASS"`
		EnableSSL bool   `envconfig:"ENABLE_SSL"`
	}

	MongoDB struct {
		Host   string `envconfig:"DB_MONGO_HOST"`
		Port   int    `envconfig:"DB_MONGO_PORT"`
		DBName string `envconfig:"DB_MONGO_NAME"`
		User   string `envconfig:"DB_MONGO_USER"`
		Pass   string `envconfig:"DB_MONGO_PASS"`
	}
}

func LoadConfig() (*Config, error) {
	// load default .env file, ignore the apperror
	_ = godotenv.Load()

	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("load config apperror: %v", err)
	}

	return cfg, nil
}
