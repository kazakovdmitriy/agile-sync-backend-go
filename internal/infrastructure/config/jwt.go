package config

type JWTConfig struct {
	SecretKey       string `mapstructure:"secret_key"`
	Issuer          string `mapstructure:"issuer"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
}
