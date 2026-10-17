package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port int
	Env  string
	DB   DBConfig
	JWT  JWTConfig
}

type DBConfig struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
	MaxConnLifeTime time.Duration
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

func Load() (*Config, error) {
	port, err := getEnvAsInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "shivam")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "postgres")
	dbSSL := getEnv("DB_SSLMODE", "disable")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSL,
	)

	jwtSecret := getEnv("JWT_SECRET", "super-secret-development-key-change-in-prod")
	jwtExpiryHours, err := getEnvAsInt("JWT_EXPIRY_HOURS", 24)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}

	return &Config{
		Port: port,
		Env:  getEnv("ENV", "development"),
		DB: DBConfig{
			URL:             dbURL,
			MaxConns:        int32(getEnvAsIntOrFallback("DB_MAX_CONNS", 25)),
			MinConns:        int32(getEnvAsIntOrFallback("DB_MIN_CONNS", 5)),
			MaxConnIdleTime: 15 * time.Minute,
			MaxConnLifeTime: 1 * time.Hour,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
			Expiry: time.Duration(jwtExpiryHours) * time.Hour,
		},
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) (int, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(valStr)
}

func getEnvAsIntOrFallback(key string, defaultVal int) int {
	val, err := getEnvAsInt(key, defaultVal)
	if err != nil {
		return defaultVal
	}
	return val
}
