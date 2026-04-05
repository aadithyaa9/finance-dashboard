package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBURL          string
	JWTSecret      string
	JWTExpiryHours int
}

func Load() *Config {
	_ = godotenv.Load()

	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		log.Fatal("JWT_EXPIRY_HOURS must be a number")
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DBURL:          getEnv("DB_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: expiryHours,
	}

	if cfg.DBURL == "" {
		log.Fatal("DB_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
