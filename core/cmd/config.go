package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
)

type ServerConfig struct {
	// General settings.
	Port   int64  `env:"PORT"`
	Logger string `env:"LOGGER"`
	Level  string `env:"LEVEL"`

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
	DemoMode       bool `env:"DEMO_MODE"`
}

type AppDBConfig struct {
	Host string `env:"APP_DATABASE_HOST"`
}

type AnalyticsDBConfig struct {
	Host string `env:"ANALYTICS_DATABASE_HOST"`
}

// This is a runtime config that is read from user settings in the database.
type RuntimeConfig struct {
	// Tracker settings.
	// Choose what type of script to serve from /script.js.
	//
	// Options:
	//
	// - "default" - Default script that collects page view data.
	//
	// - "tagged-events" - Script that collects page view data and custom event properties.
	ScriptType string
	// Usage settings.
	// Number of threads to use for processing events.
	Threads int
	// Memory limit for processing events.
	MemoryLimit string
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

	// Misc constants.
	DefaultDemoMode = false
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
		DemoMode:             DefaultDemoMode,
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

// NewRuntimeConfig creates a new runtime config.
func NewRuntimeConfig(ctx context.Context, user *sqlite.Client, analytics *duckdb.Client) (*RuntimeConfig, error) {
	// Load the script type from the database.
	settings, err := user.GetSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "runtime config")
	}

	// Set the DuckDB settings.
	if settings.MemoryLimit != "" || settings.Threads != 0 {
		err := analytics.SetDuckDBSettings(ctx, &settings.DuckDBSettings)
		if err != nil {
			return nil, errors.Wrap(err, "runtime config")
		}
	}

	return &RuntimeConfig{
		ScriptType:  settings.ScriptType,
		Threads:     settings.Threads,
		MemoryLimit: settings.MemoryLimit,
	}, nil
}
