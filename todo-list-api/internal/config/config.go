package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	port := getEnv("PORT", "8080")
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return &Config{
		Port:        port,
		DatabaseURL: dbUrl,
	}, nil
}

func getEnv(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return val
}
