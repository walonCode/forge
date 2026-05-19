package configs

import (
	"api/pkg/utils"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PORT        int
	DATABSE_URL string
	JWT_SECRET  string
	APP_VERSION string
}

func Load() *Config {
	database_url := utils.GetEnvString("DATABASE_URL", "")
	port := utils.GetEnvInt("PORT", 1000)
	appVersion := utils.GetEnvString("APP_VERSION", "1.0.0")

	return &Config{
		PORT:        port,
		DATABSE_URL: database_url,
		APP_VERSION: appVersion,
	}
}
