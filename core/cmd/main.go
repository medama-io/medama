package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/duckdb/duckdb-go/v2"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Build information. Populated at build time.
var (
	Version = ""
	Commit  = ""
)

const (
	// Default error exit code.
	exitCodeError = 1
	// Signifies the usage of some invalid shell built-in command.
	exitCodeInvalidShell = 2
)

func main() {
	err := run(context.Background(), os.Args[1:])

	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, "Usage: medama <command> [flags]")
		os.Exit(exitCodeInvalidShell)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(exitCodeError)
	}
}

func run(ctx context.Context, args []string) error {
	// Get command name
	var cmd string
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "start":
		// Check for --env flag to set configuration to also scan
		// for environment variables. Ignore all other flags.
		useEnv := false

		for _, arg := range args {
			switch arg {
			case "--env", "-env":
				useEnv = true
			}
		}

		// Create start command
		s, err := NewStartCommand(useEnv, Version, Commit)
		if err != nil {
			return err
		}

		// Parse flags
		if err := s.ParseFlags(args); err != nil {
			return err
		}

		return s.Run(ctx)

	case "version":
		//nolint: forbidigo // This is a CLI tool, so it's fine to use fmt for this part.
		fmt.Println(GetVersion())
		return nil

	default:
		if cmd == "" {
			return flag.ErrHelp
		}

		return fmt.Errorf("medama: unknown command %q", cmd)
	}
}

func GetVersion() string {
	if Version != "" && Commit != "" {
		return fmt.Sprintf("Medama Analytics %s, commit=%s", Version, Commit)
	}

	if Commit != "" {
		return "Medama Analytics commit=" + Commit
	}

	return "Medama Development Build"
}
