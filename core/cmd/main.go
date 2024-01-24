package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
)

// Build information.
var (
	//nolint: gochecknoglobals // These variables are populated at build time.
	Version = ""
	//nolint: gochecknoglobals // These variables are populated at build time.
	Commit = ""
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
		s, err := NewStartCommand()
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
		return fmt.Sprintf("Medama v%s, commit=%s", Version, Commit)
	}

	if Commit != "" {
		return fmt.Sprintf("Medama commit=%s", Commit)
	}

	return "Medama Development Build"
}
