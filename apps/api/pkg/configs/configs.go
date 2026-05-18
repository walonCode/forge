package configs

import (
	"os"
	"strconv"

	_"github.com/joho/godotenv/autoload"
)



type Config struct {
	PORT int
	DATABSE_URL string
	JWT_SECRET string
	APP_VERSION string
}

func Load()(*Config){
	database_url := os.Getenv("DATABASE_URL")
	port,_ := strconv.Atoi(os.Getenv("PORT"))
	appVersion := os.Getenv("APP_VERSION")
	
	return &Config{
		PORT: port,
		DATABSE_URL: database_url,
		APP_VERSION: appVersion,
	}
}

