package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
)

func ErrorHandler() func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		var (
			code    = http.StatusInternalServerError
			ogenErr ogenerrors.Error
		)
		switch {
		case errors.Is(err, ht.ErrNotImplemented):
			code = http.StatusNotImplemented
		case errors.As(err, &ogenErr):
			code = ogenErr.Code()
		}

		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte(fmt.Sprintf("{\"status\": %d, \"message\": \"Not Found\"}", code)))
		if err != nil {
			slog.Error("Unable to write error handler response", err)
		}

		attributes := []slog.Attr{
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Int("status_code", code),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(r.Context(), slog.LevelError, "Error", attributes...)
	}
}
