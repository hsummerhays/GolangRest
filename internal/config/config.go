package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	RateLimit      float64
	RateLimitBurst int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using defaults")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "products.db"
	}

	rateLimit := 2.0
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			rateLimit = parsed
		}
	}

	rateLimitBurst := 5
	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			rateLimitBurst = parsed
		}
	}

	return &Config{
		Port:           port,
		DBPath:         dbPath,
		RateLimit:      rateLimit,
		RateLimitBurst: rateLimitBurst,
	}
}
