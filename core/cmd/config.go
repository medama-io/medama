package main

import (
	"os"
	"time"

	"github.com/medama-io/medama/util"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	AppEnv string `env:"APP_ENV,default=development"`
	Debug  bool   `env:"DEBUG,default=false"`
	Port   uint16 `env:"PORT,default=8080"`

	// Cache settings
	CacheCleanupInterval time.Duration

	// Timeout settings
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration
	TimeoutIdle  time.Duration
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
)

// NewConfig creates a new config instance with default values.
func NewConfig() *Config {
	isDebug := os.Getenv("DEBUG") == "true"
	util.NewLogger(os.Stdout, isDebug)

	return &Config{
		Server: ServerConfig{
			AppEnv:               AppEnvDevelopment,
			Debug:                isDebug,
			Port:                 DefaultPort,
			CacheCleanupInterval: DefaultCacheCleanupInterval,
			TimeoutRead:          DefaultTimeoutRead,
			TimeoutWrite:         DefaultTimeoutWrite,
			TimeoutIdle:          DefaultTimeoutIdle,
		},
	}
}
