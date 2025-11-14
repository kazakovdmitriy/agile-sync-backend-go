package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr         string
	DatabaseURL        string
	Environment        string
	LogLevel           string
	RedisURL           string
	MaxConnections     int
	ReadTimeout        int // в секундах
	WriteTimeout       int // в секундах
	SessionTimeout     int // в минутах
	RateLimit          int // запросов в минуту
	DBMaxOpenConns     int
	DBMaxIdleConns     int
	DBConnMaxLifetime  int // в минутах
	JWTSecret          string
	JWTIssuer          string
	JWTAccessTokenTTL  int
	JWTRefreshTokenTTL int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	config := &Config{
		ServerAddr:         getEnv("SERVER_PORT", ":8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RedisURL:           getEnv("REDIS_URL", ""),
		MaxConnections:     getEnvAsInt("MAX_CONNECTIONS", 100),
		ReadTimeout:        getEnvAsInt("READ_TIMEOUT", 10),
		WriteTimeout:       getEnvAsInt("WRITE_TIMEOUT", 10),
		SessionTimeout:     getEnvAsInt("SESSION_TIMEOUT", 30),
		RateLimit:          getEnvAsInt("RATE_LIMIT", 1000),
		DBConnMaxLifetime:  getEnvAsInt("DB_CONN_MAX_LIFETIME", 25),
		DBMaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		DBMaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", 5),
		JWTSecret:          getEnv("SECRET_KEY", "default-secret-key"),
		JWTIssuer:          getEnv("JWT_ISSUER", "your-app-backend"),
		JWTAccessTokenTTL:  getEnvAsInt("ACCESS_TOKEN_EXPIRE_MINUTES", 30), // 30 минут по умолчанию
		JWTRefreshTokenTTL: getEnvAsInt("REFRESH_TOKEN_EXPIRE_DAYS", 7),    // 7 дней по умолчанию
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid value for %s: %s", key, valueStr))
	}

	return value
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) GetAccessTokenTTL() time.Duration {
	return time.Duration(c.JWTAccessTokenTTL) * time.Minute
}

func (c *Config) GetRefreshTokenTTL() time.Duration {
	return time.Duration(c.JWTRefreshTokenTTL) * time.Hour * 24
}
