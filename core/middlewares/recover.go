package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/ogen-go/ogen/middleware"
)

// Recovery is a middleware that recovers from panics.
func Recovery() middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		recovered := false

		defer func() {
			if rvr := recover(); rvr != nil {
				if errors.Is(rvr.(error), http.ErrAbortHandler) {
					panic(rvr)
				}

				slog.Log(req.Context, slog.LevelError, "panic recovery error", "error", rvr, "stack", string(debug.Stack()))

				req.Raw.Header.Add("Connection", "close")

				recovered = true
			}
		}()

		if recovered {
			return middleware.Response{}, errors.New("the server encountered a problem and could not process your request")
		}

		return next(req)
	}
}
