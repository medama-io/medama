package util

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// SetupLogger sets the default logger.
func SetupLogger(ctx context.Context, logger string, level string) (context.Context, error) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	switch logger {
	case "json":
		// Do nothing
	case "pretty":
		log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	default:
		return nil, fmt.Errorf("invalid logger type \"%s\"", logger)
	}

	switch level {
	case "debug":
		log = log.Level(zerolog.DebugLevel)
		log.Debug().Msg("Logging level set to debug")
	case "info":
		log = log.Level(zerolog.InfoLevel)
	case "warn":
		log = log.Level(zerolog.WarnLevel)
		log.Info().Msg("Logging level set to warn")
	case "error":
		log = log.Level(zerolog.ErrorLevel)
		log.Info().Msg("Logging level set to error")
	default:
		return nil, fmt.Errorf("invalid log level \"%s\"", level)
	}

	return log.WithContext(ctx), nil
}
