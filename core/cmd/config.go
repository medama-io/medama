package main

import (
	"time"

	"github.com/caarlos0/env/v10"
)

type ServerConfig struct {
	AppEnv string `env:"APP_ENV"`
	Port   int64  `env:"PORT"`

	// Cache settings
	CacheCleanupInterval time.Duration

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
	// App Environments.
	AppEnvDevelopment = "development"
	AppEnvProduction  = "production"
	DefaultPort       = 8080

	// Cache constants.
	DefaultCacheCleanupInterval = 5 * time.Minute

	// HTTP server constants.
	DefaultTimeoutRead  = 5 * time.Second
	DefaultTimeoutWrite = 10 * time.Second
	DefaultTimeoutIdle  = 15 * time.Second

	// Database constants.
	DefaultSQLiteHost = "./sqlite.dev.db"
	DefaultDuckDBHost = "./duckdb.dev.db"
)

// NewServerConfig creates a new server config.
func NewServerConfig() (*ServerConfig, error) {
	config := &ServerConfig{
		AppEnv:               AppEnvDevelopment,
		Port:                 DefaultPort,
		CacheCleanupInterval: DefaultCacheCleanupInterval,
		TimeoutRead:          DefaultTimeoutRead,
		TimeoutWrite:         DefaultTimeoutWrite,
		TimeoutIdle:          DefaultTimeoutIdle,
	}

	// Load config from environment variables.
	if err := env.Parse(config); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return config, nil
}
