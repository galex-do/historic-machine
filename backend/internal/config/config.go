package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	DatabaseURL string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "password"),
			DBName:      getEnv("DB_NAME", "historical_events"),
			SSLMode:     getEnv("DB_SSL_MODE", "disable"),
			DatabaseURL: getEnv("DATABASE_URL", ""),
		},
	}
}

// getEnv returns environment variable value or default if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}