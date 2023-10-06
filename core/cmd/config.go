package main

import (
	"time"
)

type ServerConfig struct {
	AppEnv string `env:"APP_ENV,default=development"`
	Port   uint   `env:"PORT,default=8080"`

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
