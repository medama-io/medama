package middlewares

import (
	"net/http"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/util/logger"
	"github.com/ogen-go/ogen/middleware"
)

// getCode returns the http status code from the error type.
func getCode(tresp any) int {
	if _, ok := tresp.(*api.BadRequestError); ok {
		return http.StatusBadRequest
	}

	if _, ok := tresp.(*api.ConflictError); ok {
		return http.StatusConflict
	}

	if _, ok := tresp.(*api.InternalServerError); ok {
		return http.StatusInternalServerError
	}

	if _, ok := tresp.(*api.NotFoundError); ok {
		return http.StatusNotFound
	}

	if _, ok := tresp.(*api.UnauthorisedError); ok {
		return http.StatusUnauthorized
	}

	return 0
}

// RequestLogger is a middleware that logs requests.
func RequestLogger() middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		// Add the logger to request context
		log := logger.Get()

		startTime := time.Now()
		resp, err := next(req)
		duration := time.Since(startTime)

		if err == nil {
			log = log.With().
				Str("operation", req.OperationName).
				Str("operationId", req.OperationID).
				Str("method", req.Raw.Method).
				Str("path", req.Raw.URL.Path).
				Dur("duration", duration).
				Logger()

			msg := "success"

			code := getCode(resp.Type)
			if code != 0 {
				log = log.With().Int("status_code", code).Logger()

				switch code {
				case http.StatusOK:
					msg = "200 OK"

				case http.StatusCreated:
					msg = "201 created"

				case http.StatusBadRequest:
					msg = "400 bad request"

				case http.StatusUnauthorized:
					msg = "401 unauthorised"

				case http.StatusNotFound:
					msg = "404 not found"

				case http.StatusConflict:
					msg = "409 conflict"

				case http.StatusInternalServerError:
					msg = "500 internal server error"
					log.Error().Err(err).Msg(msg)

					return resp, err
				}
			}

			log.Info().Msg(msg)
		}

		return resp, err
	}
}
