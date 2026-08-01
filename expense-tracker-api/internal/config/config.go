package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecretKey string
}

func LoadEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("missing .env file")
	}
	port := getEnv("PORT", "8080")
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if jwtSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	return &Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		JWTSecretKey: jwtSecretKey,
	}, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
