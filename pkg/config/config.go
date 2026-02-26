package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port               string
	Env                string
	JWTSecret          string
	CORSOrigins        []string
	DatabaseURL        string
	SuperAdminEmail    string
	SuperAdminPassword string
	SuperAdminName     string
}

// loadEnv reads .env file and sets environment variables. Ignores if file does not exist.
func loadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			// Remove surrounding quotes if present
			if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
			os.Setenv(key, val)
		}
	}
}

func Load() (Config, error) {
	loadEnv()

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

	corsOrigins := parseCORSOrigins(os.Getenv("CORS_ORIGINS"), env)

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("missing required environment variable: DATABASE_URL")
	}

	return Config{
		Port:               port,
		Env:                env,
		JWTSecret:          jwtSecret,
		CORSOrigins:        corsOrigins,
		DatabaseURL:        databaseURL,
		SuperAdminEmail:    os.Getenv("SUPERADMIN_EMAIL"),
		SuperAdminPassword: os.Getenv("SUPERADMIN_PASSWORD"),
		SuperAdminName:     os.Getenv("SUPERADMIN_NAME"),
	}, nil
}

func parseCORSOrigins(envVal, goEnv string) []string {
	if envVal != "" {
		origins := strings.Split(envVal, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		return origins
	}
	if goEnv == "development" {
		return []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:8081", "http://127.0.0.1:3000", "http://127.0.0.1:8080", "http://127.0.0.1:8081"}
	}
	return nil
}
