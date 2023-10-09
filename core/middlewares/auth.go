package middlewares

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
)

type Handler struct {
	cache *util.Cache
}

// Compile time check for Handler.
var _ api.SecurityHandler = (*Handler)(nil)

// NewAuthHandler returns a new instance of the auth service handler.
func NewAuthHandler(cache *util.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

// HandleCookieAuth handles cookie based authentication.
func (h *Handler) HandleCookieAuth(ctx context.Context, _operationName string, t api.CookieAuth) (context.Context, error) {
	// Check if session exists
	_, err := h.cache.Get(ctx, t.APIKey)
	// If session does not exist, return error
	if err != nil {
		return nil, model.ErrUnauthorised
	}

	return ctx, nil
}
