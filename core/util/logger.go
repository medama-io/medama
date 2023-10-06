package util

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/ogen-go/ogen/middleware"
)

// NewHandler creates a new slog handler.
func NewHandler(w io.Writer, isDebug bool) slog.Handler {
	logLevel := &slog.LevelVar{}

	handlerOpts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	switch getAppEnv() {
	case "production":
		handler = slog.NewJSONHandler(w, handlerOpts)
	case "development":
		// Use tint handler for pretty readable logs
		tintOpts := &tint.Options{
			Level: logLevel,
		}

		if w == os.Stdout {
			handler = tint.NewHandler(colorable.NewColorableStdout(), tintOpts)
		} else {
			handler = slog.NewTextHandler(w, handlerOpts)
		}
	default:
		handler = slog.NewTextHandler(w, handlerOpts)
	}

	if isDebug {
		logLevel.Set(slog.LevelDebug)
	}

	return handler
}

// SetupLogger sets the default logger.
func SetupLogger(w io.Writer, isDebug bool) {
	handler := NewHandler(w, isDebug)
	slog.SetDefault(slog.New(handler))

	if isDebug {
		slog.Debug("Debug logging enabled")
	}
}

func RequestLogger() middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		resp, err := next(req)

		if err != nil {
			slog.Error("Error", err)
		} else {
			attributes := []slog.Attr{
				slog.String("operation", req.OperationName),
				slog.String("operationId", req.OperationID),
				slog.String("method", req.Raw.Method),
				slog.String("path", req.Raw.URL.Path),
			}

			if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
				attributes = append(attributes, slog.Int("status_code", tresp.GetStatusCode()))
			}

			slog.LogAttrs(req.Context, slog.LevelInfo, "Success", attributes...)
		}
		return resp, err
	}
}

// Helper function to get app env from env var before initializing config used in logger.
func getAppEnv() string {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return "development"
	}
	return appEnv
}
