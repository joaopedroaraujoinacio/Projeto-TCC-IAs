package config

import (
	"os"
	"path/filepath"
)


type Config struct {
	DatabaseURL string
	Port 				string
}

func Load() *Config {
	return &Config{
 		DatabaseURL: getEnv("DATABASE_URL", getDefaultDatabasePath()),
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDefaultDatabasePath() string {
	return filepath.Join("data", "rag_database.db")
}

