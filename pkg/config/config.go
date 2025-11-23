package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Env       string
	JWTSecret string
}

func Load() (Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("missing required environment variable: PORT")
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		return Config{}, fmt.Errorf("missing required environment variable: GO_ENV")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("missing required environment variable: JWT_SECRET")
	}

	return Config{
		Port:      port,
		Env:       env,
		JWTSecret: jwtSecret,
	}, nil
}
