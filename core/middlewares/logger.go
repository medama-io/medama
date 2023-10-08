package middlewares

import (
	"log/slog"

	"github.com/ogen-go/ogen/middleware"
)

// RequestLogger is a middleware that logs requests.
func RequestLogger() middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		resp, err := next(req)

		if err != nil {
			slog.ErrorContext(req.Context, "Error", err)
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
