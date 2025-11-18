package config

type ServerConfig struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"read_timeout"`
	WriteTimeout   int    `mapstructure:"write_timeout"`
	MaxConnections int    `mapstructure:"max_connections"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
}
