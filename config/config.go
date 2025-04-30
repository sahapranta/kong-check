package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSL            string
	DefaultProtocol  string
	DefaultHostname  string
	RequestTimeout   int
	MaxConcurrentReq int
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		// Only log warning as environment variables might be set directly
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		DBHost:           getEnv("PG_HOST", "localhost"),
		DBPort:           getEnv("PG_PORT", "5432"),
		DBUser:           getEnv("PG_USER", "kong"),
		DBPassword:       getEnv("PG_PASSWORD", "kong"),
		DBName:           getEnv("PG_DATABASE", "kong"),
		DBSSL:            getEnv("PG_SSLMODE", "disable"),
		DefaultProtocol:  getEnv("DEFAULT_PROTOCOL", "http"),
		DefaultHostname:  getEnv("DEFAULT_HOSTNAME", "localhost:8000"),
		RequestTimeout:   5, // Default timeout in seconds
		MaxConcurrentReq: 10,
	}

	// Build connection string
	config.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSL,
	)

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
