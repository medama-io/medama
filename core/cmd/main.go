package main

import (
	"context"
	"flag"
	"fmt"
	"os"
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

	if err == flag.ErrHelp {
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
	case "version":
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
