package main

import (
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-faster/errors"
)

type ServerConfig struct {
	Port int64 `env:"PORT"`

	// Cache settings
	CacheCleanupInterval time.Duration

	// CORS Settings
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`

	// Logging settings
	Logger string `env:"LOGGER"`
	Level  string `env:"LOGGER_LEVEL"`

	// Timeout settings
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration
	TimeoutIdle  time.Duration
}

type AppDBConfig struct {
	Host string `env:"APP_DATABASE_HOST"`
}

type AnalyticsDBConfig struct {
	Host string `env:"ANALYTICS_DATABASE_HOST"`
}

const (
	// General constants.
	DefaultPort = 8080

	// Cache constants.
	DefaultCacheCleanupInterval = 5 * time.Minute

	// HTTP server constants.
	DefaultTimeoutRead  = 5 * time.Second
	DefaultTimeoutWrite = 10 * time.Second
	DefaultTimeoutIdle  = 15 * time.Second

	// Database constants.
	DefaultSQLiteHost = "./sqlite.dev.db"
	DefaultDuckDBHost = "./duckdb.dev.db"

	// Logging constants.
	DefaultLogger      = "json"
	DefaultLoggerLevel = "info"
)

// NewServerConfig creates a new server config.
func NewServerConfig() (*ServerConfig, error) {
	config := &ServerConfig{
		Port:                 DefaultPort,
		CacheCleanupInterval: DefaultCacheCleanupInterval,
		Logger:               DefaultLogger,
		Level:                DefaultLoggerLevel,
		TimeoutRead:          DefaultTimeoutRead,
		TimeoutWrite:         DefaultTimeoutWrite,
		TimeoutIdle:          DefaultTimeoutIdle,
	}

	// Load config from environment variables.
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "config")
	}

	return config, nil
}

// NewAppDBConfig creates a new app database config.
func NewAppDBConfig() (*AppDBConfig, error) {
	config := &AppDBConfig{
		Host: DefaultSQLiteHost,
	}

	// Load config from environment variables.
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "config")
	}

	return config, nil
}

// NewAnalyticsDBConfig creates a new analytics database config.
func NewAnalyticsDBConfig() (*AnalyticsDBConfig, error) {
	config := &AnalyticsDBConfig{
		Host: DefaultDuckDBHost,
	}

	// Load config from environment variables.
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "config")
	}

	return config, nil
}
