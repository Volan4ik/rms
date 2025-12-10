package config

import (
	"log"
	"os"
)

// Config holds application configuration loaded from environment.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	HTTPPort   string
}

// Load reads environment variables with sensible defaults for local development.
func Load() *Config {
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "rms"),
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
	}

	if cfg.DBPassword == "" {
		log.Println("warning: DB_PASSWORD is empty, set it in environment for production")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
