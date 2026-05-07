package config

import (
	"os"
	"path/filepath"
	"github.com/joho/godotenv"
)


type Config struct {
	DatabaseURL 	string
	Port 					string
	OpenAIAPIKey 	string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
 		DatabaseURL: getEnv("DATABASE_URL", getDefaultDatabasePath()),
		Port: getEnv("PORT", "8080"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
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

