package util

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// SetupLogger sets the default logger.
func SetupLogger(ctx context.Context, isDebug bool, logger string) context.Context {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if logger == "pretty" {
		log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	if isDebug {
		log = log.Level(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	}

	return log.WithContext(ctx)
}
