package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		DBHost:     getEnvOrPanic("DB_HOST"),
		DBUser:     getEnvOrPanic("DB_USER"),
		DBPassword: getEnvOrPanic("DB_PASSWORD"),
		DBName:     getEnvOrPanic("DB_NAME"),
		DBPort:     getEnvOrPanic("DB_PORT"),
		JWTSecret:  getEnvOrPanic("JWT_SECRET"),
		ServerPort: getEnvOrPanic("SERVER_PORT"),
	}

	return config, nil
}

func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return value
}
