package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(name, default_value string) string {
	value := os.Getenv(name)
	if strings.TrimSpace(value) == "" {
		return default_value
	}

	return value
}

func GetEnvInt(name string, default_value int) int {
	value := os.Getenv(name)
	envInt, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return default_value
	}

	return envInt
}
