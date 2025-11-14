package config

import "time"

type JWTConfig struct {
	SecretKey       string        `mapstructure:"secret_key"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Issuer          string        `mapstructure:"issuer"`
}

func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:       "your-super-secret-key-change-in-production",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour, // 7 дней
		Issuer:          "your-app-name",
	}
}
