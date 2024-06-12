package middlewares

import (
	"errors"
	"net/http"

	"github.com/medama-io/medama/util/logger"
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

				log := logger.Get()
				log.Error().
					Str("path", req.Raw.URL.Path).
					Str("method", req.Raw.Method).
					Str("User-Agent", req.Raw.Header.Get("User-Agent")).
					Stack().
					Err(rvr.(error)).
					Msg("panic recovery error")

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
