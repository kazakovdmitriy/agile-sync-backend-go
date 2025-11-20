package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type Config struct {
	Environment        string          `mapstructure:"environment"`
	LogLevel           string          `mapstructure:"log_level"`
	GuestCleanInterval time.Duration   `mapstructure:"guest_clean_interval"`
	Server             ServerConfig    `mapstructure:"server"`
	Database           DatabaseConfig  `mapstructure:"database"`
	JWT                JWTConfig       `mapstructure:"jwt"`
	WebSocket          WebSocketConfig `mapstructure:"websocket"`
	RateLimit          RateLimitConfig `mapstructure:"rate_limit"`
}

func LoadConfig(configPath string) (*Config, error) {

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	setDefaults()

	bindEnvVars()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // server.addr → SERVER_ADDR

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	processEnvVars()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

func processEnvVars() {
	// Обрабатываем WEBSOCKET_ALLOWED_ORIGINS и CORS_ORIGINS
	processStringSliceEnvVar("CORS_ORIGINS", "websocket.allowed_origins")
}

func processStringSliceEnvVar(envVar string, configPath string) {
	if envValue := viper.GetString(envVar); envValue != "" {
		values := strings.Split(envValue, ",")
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
		viper.Set(configPath, values)
	}
}

// Вспомогательные методы
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetServerReadTimeout возвращает timeout в виде time.Duration
func (c *Config) GetServerReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeout) * time.Second
}

// GetServerWriteTimeout возвращает timeout в виде time.Duration
func (c *Config) GetServerWriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeout) * time.Second
}

// GetDBConnMaxLifetime возвращает lifetime в виде time.Duration
func (c *Config) GetDBConnMaxLifetime() time.Duration {
	return time.Duration(c.Database.ConnMaxLifetime) * time.Minute
}

func (c *Config) GetAccessTokenTTL() time.Duration {
	return time.Duration(c.JWT.AccessTokenTTL) * time.Minute
}

func (c *Config) GetRefreshTokenTTL() time.Duration {
	return time.Duration(c.JWT.RefreshTokenTTL) * 24 * time.Hour
}

// IsOriginAllowed проверяет, разрешен ли origin для WebSocket соединений
func (c *Config) IsOriginAllowed(origin string) bool {
	if len(c.WebSocket.AllowedOrigins) == 0 {
		return false
	}

	for _, allowed := range c.WebSocket.AllowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
		if strings.Contains(allowed, "*") {
			if matchWildcard(allowed, origin) {
				return true
			}
		}
	}
	return false
}

func matchWildcard(pattern, origin string) bool {
	patternParts := strings.Split(pattern, "*")
	if len(patternParts) != 2 {
		return false
	}

	return strings.HasPrefix(origin, patternParts[0]) &&
		strings.HasSuffix(origin, patternParts[1])
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.max_connections", 100)

	// Database defaults
	viper.SetDefault("database.url", "postgres://user:password@localhost:5432/dbname?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 5)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", 25)

	// JWT defaults
	viper.SetDefault("jwt.secret_key", "default-secret-key")
	viper.SetDefault("jwt.issuer", "your-app-backend")
	viper.SetDefault("jwt.access_token_ttl", 30*time.Minute)
	viper.SetDefault("jwt.refresh_token_ttl", 7*24*time.Hour)

	// WebSocket defaults - согласно вашим параметрам
	viper.SetDefault("websocket.enabled", true)
	viper.SetDefault("websocket.allowed_origins", []string{"*"})
	viper.SetDefault("websocket.handshake_timeout", 10*time.Second)
	viper.SetDefault("websocket.write_timeout", 10*time.Second)
	viper.SetDefault("websocket.read_timeout", 60*time.Second)
	viper.SetDefault("websocket.pong_timeout", 60*time.Second)
	viper.SetDefault("websocket.ping_interval", 30*time.Second)
	viper.SetDefault("websocket.read_buffer_size", 4096)
	viper.SetDefault("websocket.write_buffer_size", 4096)
	viper.SetDefault("websocket.max_message_size", 1048576) // 1MB
	viper.SetDefault("websocket.max_connections_per_session", 100)
	viper.SetDefault("websocket.enable_compression", true)

	// Other defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("session.timeout", 30)
	viper.SetDefault("rate_limit.requests_per_minute", 1000)
	viper.SetDefault("guest_clean_interval", time.Duration(12)*time.Hour)
}

func bindEnvVars() {
	// Server env vars
	viper.BindEnv("server.addr", "SERVER_ADDR")
	viper.BindEnv("server.read_timeout", "READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "WRITE_TIMEOUT")
	viper.BindEnv("server.max_connections", "MAX_CONNECTIONS")

	// Database env vars
	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")

	// Redis env vars
	viper.BindEnv("redis.url", "REDIS_URL")
	viper.BindEnv("websocket.handshake_timeout", "WEBSOCKET_HANDSHAKE_TIMEOUT")

	// JWT env vars
	viper.BindEnv("jwt.secret_key", "JWT_SECRET")
	viper.BindEnv("jwt.issuer", "JWT_ISSUER")
	viper.BindEnv("jwt.access_token_ttl", "ACCESS_TOKEN_EXPIRE_MINUTES")
	viper.BindEnv("jwt.refresh_token_ttl", "REFRESH_TOKEN_EXPIRE_DAYS")

	// Websoket
	viper.BindEnv("websocket.allowed_origins", "CORS_ORIGINS")
	viper.BindEnv("websocket.enabled", "WEBSOCKET_ENABLED")
	viper.BindEnv("websocket.handshake_timeout", "WEBSOCKET_HANDSHAKE_TIMEOUT")
	viper.BindEnv("websocket.write_timeout", "WEBSOCKET_WRITE_TIMEOUT")
	viper.BindEnv("websocket.read_timeout", "WEBSOCKET_READ_TIMEOUT")
	viper.BindEnv("websocket.pong_timeout", "WEBSOCKET_PONG_TIMEOUT")
	viper.BindEnv("websocket.ping_interval", "WEBSOCKET_PING_INTERVAL")
	viper.BindEnv("websocket.read_buffer_size", "WEBSOCKET_READ_BUFFER_SIZE")
	viper.BindEnv("websocket.write_buffer_size", "WEBSOCKET_WRITE_BUFFER_SIZE")
	viper.BindEnv("websocket.max_message_size", "WEBSOCKET_MAX_MESSAGE_SIZE")
	viper.BindEnv("websocket.max_connections_per_session", "WEBSOCKET_MAX_CONNECTIONS_PER_SESSION")
	viper.BindEnv("websocket.enable_compression", "WEBSOCKET_ENABLE_COMPRESSION")

	// Other env vars
	viper.BindEnv("environment", "ENVIRONMENT")
	viper.BindEnv("log_level", "LOG_LEVEL")
	viper.BindEnv("guest_clean_interval", "GUEST_CLEAN_INTERVAL")
	viper.BindEnv("session.timeout", "SESSION_TIMEOUT")
	viper.BindEnv("rate_limit.requests_per_minute", "RATE_LIMIT")
}
