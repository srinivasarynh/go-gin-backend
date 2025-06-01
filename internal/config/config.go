package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port string
	DatabaseURL string
	RedisURL string
	JWTSecret string
	RateLimitRPS int
	LogLevel string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port: getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "my-super-secret-key"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		RateLimitRPS: rateLimitRPS,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
