package main

import (
	"context"
	"flag"
)

type StartCommand struct {
	// Enable verbose debug logging
	IsDebug bool
}

// NewStartCommand creates a new start command.
func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

// ParseFlags parses the command line flags for the start command.
func (s *StartCommand) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.BoolVar(&s.IsDebug, "debug", false, "Enable verbose debug logging")

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	return nil
}

// Run executes the start command.
func (s *StartCommand) Run(ctx context.Context) error {
	// Create new config
	c := NewConfig()
	c.Server.Debug = s.IsDebug

	return nil
}
