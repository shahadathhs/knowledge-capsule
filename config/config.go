package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
	Env  string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("missing required environment variable: PORT")
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		return Config{}, fmt.Errorf("missing required environment variable: GO_ENV")
	}

	return Config{
		Port: port,
		Env:  env,
	}, nil
}
