package main

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-faster/errors"
)

type ServerConfig struct {
	// General settings.
	Port   int64  `env:"PORT"`
	Logger string `env:"LOGGER"`
	Level  string `env:"LEVEL"`

	// AutoSSL settings.
	AutoSSLDomain string `env:"AUTO_SSL"`
	AutoSSLEmail  string `env:"AUTO_SSL_EMAIL"`

	// Cache settings.
	CacheCleanupInterval time.Duration

	// CORS Settings.
	//nolint: tagalign // It removes the comma.
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`

	// Timeout settings.
	TimeoutReadHeader time.Duration
	TimeoutRead       time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration

	// Misc settings.
	// Enable /debug/pprof endpoints.
	Profiler       bool `env:"PROFILER"`
	UseEnvironment bool
	DemoMode       bool `env:"DEMO_MODE"`

	// Short Git commit SHA. Used when returning the version of the server in
	// X-API-Commit header for client-side cache busting.
	Commit  string
	Version string
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
	DefaultTimeoutReadHeader = 10 * time.Second
	DefaultTimeoutRead       = 30 * time.Second
	DefaultTimeoutWrite      = 2 * time.Minute
	DefaultTimeoutIdle       = 5 * time.Minute

	// Database constants.
	DefaultSQLiteHost = "./me_meta.db"
	DefaultDuckDBHost = "./me_analytics.db"

	// Logging constants.
	DefaultLogger      = "json"
	DefaultLoggerLevel = "info"

	// Misc constants.
	DefaultProfiler = false
	DefaultDemoMode = false
)

// NewServerConfig creates a new server config.
func NewServerConfig(useEnv bool, version string, commit string) (*ServerConfig, error) {
	if version == "" {
		version = "development"
	}

	if commit == "" {
		commit = "development"
	}

	config := &ServerConfig{
		Port:                 DefaultPort,
		CacheCleanupInterval: DefaultCacheCleanupInterval,
		Logger:               DefaultLogger,
		Level:                DefaultLoggerLevel,
		TimeoutReadHeader:    DefaultTimeoutReadHeader,
		TimeoutRead:          DefaultTimeoutRead,
		TimeoutWrite:         DefaultTimeoutWrite,
		TimeoutIdle:          DefaultTimeoutIdle,
		Profiler:             DefaultProfiler,
		UseEnvironment:       useEnv,
		DemoMode:             DefaultDemoMode,
		Version:              version,
		Commit:               commit,
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
