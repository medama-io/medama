package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-faster/jx"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"github.com/ogen-go/ogen/ogenerrors"
)

// ErrorHandler is a middleware that handles any unhandled errors by ogen.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	errMessage := strings.ReplaceAll(err.Error(), "\"", "'")

	log := logger.Get().With().
		Str("path", req.URL.Path).
		Str("method", req.Method).
		Int("status_code", code).
		Str("message", errMessage).
		Str("Connection", req.Header.Get("Connection")).
		Str("Content-Type", req.Header.Get("Content-Type")).
		Str("Content-Length", req.Header.Get("Content-Length")).
		Str("User-Agent", req.Header.Get("User-Agent")).Logger()

	if errors.Is(err, model.ErrUnauthorised) || code == http.StatusUnauthorized {
		log.Warn().Msg("unauthorised")
	} else {
		log.Error().Msg(err.Error())
	}

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
