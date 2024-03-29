package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
)

// ErrorHandler is a middleware that handles any unhandled errors by ogen.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	errMessage := strings.ReplaceAll(err.Error(), "\"", "'")

	attributes := []slog.Attr{
		slog.String("path", req.URL.Path),
		slog.String("method", req.Method),
		slog.Int("status_code", code),
		slog.String("message", errMessage),
		slog.String("Connection", req.Header.Get("Connection")),
		slog.String("Content-Type", req.Header.Get("Content-Type")),
		slog.String("Content-Length", req.Header.Get("Content-Length")),
		slog.String("User-Agent", req.Header.Get("User-Agent")),
	}
	slog.LogAttrs(ctx, slog.LevelError, "error", attributes...)

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		errMessage = "missing security token or cookie"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.ObjStart()
	e.FieldStart("code")
	e.Int(code)
	e.FieldStart("message")
	e.StrEscape(errMessage)
	e.ObjEnd()
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
