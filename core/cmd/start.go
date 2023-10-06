package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/medama-io/medama/util"
)

type StartCommand struct {
	Debug  bool
	Server ServerConfig
}

// NewStartCommand creates a new start command.
func NewStartCommand() *StartCommand {
	return &StartCommand{
		Server: ServerConfig{
			AppEnv:               AppEnvDevelopment,
			CacheCleanupInterval: DefaultCacheCleanupInterval,
			TimeoutRead:          DefaultTimeoutRead,
			TimeoutWrite:         DefaultTimeoutWrite,
			TimeoutIdle:          DefaultTimeoutIdle,
		},
	}
}

// ParseFlags parses the command line flags for the start command.
func (s *StartCommand) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.BoolVar(&s.Debug, "debug", false, "Enable verbose debug logging")
	fs.UintVar(&s.Server.Port, "port", DefaultPort, "Port to listen on")

	// Parse flags
	return fs.Parse(args)
}

// Run executes the start command.
func (s *StartCommand) Run(ctx context.Context) error {
	util.SetupLogger(os.Stdout, s.Debug)
	slog.Info(GetVersion())

	return nil
}
