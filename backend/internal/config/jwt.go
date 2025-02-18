package config

import (
	"os"
	"strconv"
)

var (
	JWTSecret      = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key-here"))
	JWTExpiryHours = getEnvIntOrDefault("JWT_EXPIRY_HOURS", 24)
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}
