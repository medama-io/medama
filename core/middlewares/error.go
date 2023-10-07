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

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	errMessage := strings.ReplaceAll(errors.Unwrap(err).Error(), "\"", "'")

	attributes := []slog.Attr{
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
		slog.Int("status_code", code),
		slog.String("message", errMessage),
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
	e.StrEscape(errMessage)
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
