package utils

import (
	"os"
	"strconv"
)

// GetEnvOrDefault retrieves the value of the environment variable named by the key.
// If the variable is not present, returns the default value.
func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvIntOrDefault retrieves the integer value of the environment variable.
// If the variable is not present or cannot be converted to int, returns the default value.
func GetEnvIntOrDefault(key string, defaultValue int) int {
	strValue := GetEnvOrDefault(key, "")
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetEnvBoolOrDefault retrieves the boolean value of the environment variable.
// If the variable is not present or cannot be converted to bool, returns the default value.
func GetEnvBoolOrDefault(key string, defaultValue bool) bool {
	strValue := GetEnvOrDefault(key, "")
	if strValue == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
