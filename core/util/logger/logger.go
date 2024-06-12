package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	//nolint:gochecknoglobals // Global logger is required for the singleton pattern.
	log zerolog.Logger
	//nolint:gochecknoglobals // Global once is required for the singleton pattern.
	once sync.Once
	err  error
)

func Init(logger string, level string) (zerolog.Logger, error) {
	once.Do(func() {
		// Set the stack marshalling and time format for zerolog.
		//nolint:reassign // Reassigning is the recommended way to set stack marshalling.
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log = zerolog.New(os.Stderr).With().Timestamp().Logger()

		// Configure the logger format.
		switch logger {
		case "json":
			// Default JSON format, do nothing
		case "pretty":
			log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		default:
			err = fmt.Errorf("invalid logger type \"%s\"", logger)
			return
		}

		// Configure the log level.
		switch level {
		case "debug":
			log = log.Level(zerolog.DebugLevel)
			log.Debug().Msg("Logging level set to debug")
		case "info":
			log = log.Level(zerolog.InfoLevel)
		case "warn":
			log = log.Level(zerolog.WarnLevel)
			log.Warn().Msg("Logging level set to warn")
		case "error":
			log = log.Level(zerolog.ErrorLevel)
			log.Error().Msg("Logging level set to error")
		default:
			err = fmt.Errorf("invalid log level \"%s\"", level)
			return
		}
	})

	return log, err
}

// Get initialises and returns a singleton zerolog.Logger instance.
// logger: The logger format, either "json" or "pretty".
// level: The log level, can be "debug", "info", "warn", or "error".
func Get() zerolog.Logger {
	return log
}
