package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment       string
	Port              string
	DBUrl             string
	RedisUrl          string
	RedisPassword     string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxConnLifetime time.Duration
	DBMaxConnIdleTime time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "development")

	port := getEnv("PORT", "8080")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5435")
	dbUser := getEnv("DB_USER", "admin")
	dbPass := getEnv("DB_PASSWORD", "secretpassword")
	dbName := getEnv("DB_NAME", "fleetcommander_db")

	sslMode := "disable"
	if appEnv == "production" || appEnv == "staging" {
		sslMode = "require"
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode)

	maxConns := getEnvAsInt("DB_MAX_CONNS", 25)
	minConns := getEnvAsInt("DB_MIN_CONNS", 5)
	maxConnLifetime := getEnvAsDuration("DB_MAX_CONN_LIFETIME", 1*time.Hour)
	maxConnIdleTime := getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute)

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6385")
	redisUrl := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisPass := getEnv("REDIS_PASSWORD", "")

	return &Config{
		Environment:       appEnv,
		Port:              port,
		DBUrl:             dbUrl,
		RedisUrl:          redisUrl,
		RedisPassword:     redisPass,
		DBMaxConns:        int32(maxConns),
		DBMinConns:        int32(minConns),
		DBMaxConnLifetime: maxConnLifetime,
		DBMaxConnIdleTime: maxConnIdleTime,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return fallback
}
