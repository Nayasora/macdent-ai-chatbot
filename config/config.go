package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	ApiServer ApiServerConfig
	Database  DatabaseConfig
	Qdrant    QdrantConfig
}

type ApiServerConfig struct {
	Port string
}

type QdrantConfig struct {
	Host   string
	Port   int
	ApiKey string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewConfig() *Config {
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	qdrantPort, _ := strconv.Atoi(getEnv("QDRANT_PORT", "6334"))

	return &Config{
		ApiServer: ApiServerConfig{
			Port: getEnv("API_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "macdent_ai"),
		},
		Qdrant: QdrantConfig{
			Host:   getEnv("QDRANT_HOST", "localhost"),
			Port:   qdrantPort,
			ApiKey: getEnv("QDRANT_API_KEY", "macdent-ai-api-key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
