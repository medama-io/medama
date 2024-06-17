package main

import (
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-faster/errors"
)

type ServerConfig struct {
	// General settings.
	Port   int64  `env:"PORT"`
	Logger string `env:"LOGGER"`
	Level  string `env:"LOGGER_LEVEL"`

	// Cache settings.
	CacheCleanupInterval time.Duration

	// CORS Settings.
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`

	// Timeout settings.
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration
	TimeoutIdle  time.Duration

	// Misc settings.
	UseEnvironment bool
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
	DefaultSQLiteHost = "./me_meta.db"
	DefaultDuckDBHost = "./me_analytics.db"

	// Logging constants.
	DefaultLogger      = "json"
	DefaultLoggerLevel = "info"
)

// NewServerConfig creates a new server config.
func NewServerConfig(useEnv bool) (*ServerConfig, error) {
	config := &ServerConfig{
		Port:                 DefaultPort,
		CacheCleanupInterval: DefaultCacheCleanupInterval,
		Logger:               DefaultLogger,
		Level:                DefaultLoggerLevel,
		TimeoutRead:          DefaultTimeoutRead,
		TimeoutWrite:         DefaultTimeoutWrite,
		TimeoutIdle:          DefaultTimeoutIdle,
		UseEnvironment:       useEnv,
	}

	// Load config from environment variables.
	if useEnv {
		if err := env.Parse(config); err != nil {
			return nil, errors.Wrap(err, "config")
		}
	}

	return config, nil
}

// NewAppDBConfig creates a new app database config.
func NewAppDBConfig(useEnv bool) (*AppDBConfig, error) {
	config := &AppDBConfig{
		Host: DefaultSQLiteHost,
	}

	// Load config from environment variables.
	if useEnv {
		if err := env.Parse(config); err != nil {
			return nil, errors.Wrap(err, "config")
		}
	}

	return config, nil
}

// NewAnalyticsDBConfig creates a new analytics database config.
func NewAnalyticsDBConfig(useEnv bool) (*AnalyticsDBConfig, error) {
	config := &AnalyticsDBConfig{
		Host: DefaultDuckDBHost,
	}

	// Load config from environment variables.
	if useEnv {
		if err := env.Parse(config); err != nil {
			return nil, errors.Wrap(err, "config")
		}
	}

	return config, nil
}
