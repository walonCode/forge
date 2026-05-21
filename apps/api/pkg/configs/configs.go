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
	jwtSecret := utils.GetEnvString("JWT_SECRET", "e8d5cbe9c3a1c605d61d21908054ebb54d095d3ee28c5a4c7b4ac1620dc6fd119836e85e54f9a5a2e9f70cd1e42fb30bebefb4d1096041a475555723b191fd5a")

	return &Config{
		PORT:        port,
		DATABSE_URL: database_url,
		APP_VERSION: appVersion,
		JWT_SECRET: jwtSecret,
	}
}
