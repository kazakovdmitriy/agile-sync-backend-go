package config

import "time"

type WebSocketConfig struct {
	Enabled                  bool          `mapstructure:"enabled"`
	AllowedOrigins           []string      `mapstructure:"allowed_origins"`
	HandshakeTimeout         time.Duration `mapstructure:"handshake_timeout"`
	WriteTimeout             time.Duration `mapstructure:"write_timeout"`
	ReadTimeout              time.Duration `mapstructure:"read_timeout"`
	PongTimeout              time.Duration `mapstructure:"pong_timeout"`
	PingInterval             time.Duration `mapstructure:"ping_interval"`
	ReadBufferSize           int           `mapstructure:"read_buffer_size"`
	WriteBufferSize          int           `mapstructure:"write_buffer_size"`
	MaxMessageSize           int64         `mapstructure:"max_message_size"`
	MaxConnectionsPerSession int           `mapstructure:"max_connections_per_session"`
	EnableCompression        bool          `mapstructure:"enable_compression"`
}
