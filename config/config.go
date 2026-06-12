package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbPort := 3306
	if port := os.Getenv("DB_PORT"); port != "" {
		var err error
		dbPort, err = strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid DB_PORT: %w", err)
		}
	}

	appPort := 3000
	if port := os.Getenv("APP_PORT"); port != "" {
		var err error
		appPort, err = strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid APP_PORT: %w", err)
		}
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBName:     getEnv("DB_NAME", "usersdb"),
		AppPort:    appPort,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
