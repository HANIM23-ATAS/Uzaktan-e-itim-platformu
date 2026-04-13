// Package config centralizes application configuration via environment variables.
package config

import (
	"os"
	"strconv"
)

// Config holds all runtime configuration values.
type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	RateLimit  float64
	BurstLimit int
}

// Load reads configuration from environment variables with sane defaults.
func Load() *Config {
	rate, _ := strconv.ParseFloat(getEnv("RATE_LIMIT", "5"), 64)
	burst, _ := strconv.Atoi(getEnv("BURST_LIMIT", "10"))

	return &Config{
		Port:       getEnv("PORT", "8090"),
		DBPath:     getEnv("DB_PATH", "./golearn.db"),
		JWTSecret:  getEnv("JWT_SECRET", "golearn-super-secret-key-change-in-prod"),
		RateLimit:  rate,
		BurstLimit: burst,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
