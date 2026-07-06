package configs

import (
	"fmt"

	"api/pkg/utils"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port        int
	DatabaseURL string
	JwtSecret   string
	AppVersion  string
}

// Load reads configuration from the environment (and an optional .env file),
// returning an error when a required value is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:        utils.GetEnvInt("PORT", 8080),
		DatabaseURL: utils.GetEnvString("DATABASE_URL", ""),
		JwtSecret:   utils.GetEnvString("JWT_SECRET", ""),
		AppVersion:  utils.GetEnvString("APP_VERSION", "1.0.0"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
