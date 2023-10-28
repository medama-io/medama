package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/medama-io/medama/api"
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
		resp, err := next(req)

		if err == nil {
			attributes := []slog.Attr{
				slog.String("operation", req.OperationName),
				slog.String("operationId", req.OperationID),
				slog.String("method", req.Raw.Method),
				slog.String("path", req.Raw.URL.Path),
			}

			level := slog.LevelInfo
			msg := "success"
			code := getCode(resp.Type)
			if code != 0 {
				attributes = append(attributes, slog.Int("status_code", code))

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
					level = slog.LevelError
				}
			}

			slog.LogAttrs(req.Context, level, msg, attributes...)
		}
		return resp, err
	}
}
