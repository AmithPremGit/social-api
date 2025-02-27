package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port        string
	Env         string
	DB          DBConfig
	Redis       RedisConfig
	Auth        AuthConfig
	RateLimiter RateLimiterConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	TokenSecret   string
	TokenIssuer   string
	TokenAudience string
	TokenExpiry   time.Duration
}

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Enabled           bool
	RequestsPerWindow int
	WindowDuration    time.Duration
}

// Load loads configuration from environment variables
func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
		DB: DBConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			Database:     getEnv("DB_NAME", "socialnetwork"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  getEnvAsDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Enabled:  getEnvAsBool("REDIS_ENABLED", false),
		},
		Auth: AuthConfig{
			// TODO: set AUTH_TOKEN_SECRET environment variable in production
			TokenSecret:   getEnv("AUTH_TOKEN_SECRET", "your-super-secret-key-change-in-production"),
			TokenIssuer:   getEnv("AUTH_TOKEN_ISSUER", "social-api"),
			TokenAudience: getEnv("AUTH_TOKEN_AUDIENCE", "social-api-users"),
			TokenExpiry:   getEnvAsDuration("AUTH_TOKEN_EXPIRY", 24*time.Hour),
		},
		RateLimiter: RateLimiterConfig{
			Enabled:           getEnvAsBool("RATE_LIMITER_ENABLED", true),
			RequestsPerWindow: getEnvAsInt("RATE_LIMITER_REQUESTS", 20),
			WindowDuration:    getEnvAsDuration("RATE_LIMITER_WINDOW", 5*time.Second),
		},
	}
}

// Helper functions for environment variables

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
