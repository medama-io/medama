package middlewares

import (
	"context"

	"github.com/medama-io/medama/model"
	"github.com/ogen-go/ogen/middleware"
)

// RequestContext adds the request object to the context for use in handlers.
func RequestContext() middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		// Match with event operations
		switch req.OperationID {
		case "post-event-hit", "get-event-ping":
			// Add the request to the context
			req.Context = context.WithValue(req.Context, model.RequestKeyBody, req.Raw)
		}

		return next(req)
	}
}
