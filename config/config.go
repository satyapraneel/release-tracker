package config

import (
	"os"
	"strconv"
)

//https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
type Config struct {
	Database DatabaseConfig
}

func New() *Config {
	// godotenv.Load() //@todo check how to load only once
	return &Config{
		Database: getDBConfig(),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
