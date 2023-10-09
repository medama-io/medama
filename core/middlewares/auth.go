package middlewares

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
)

type Handler struct {
	auth *util.AuthService
}

// Compile time check for Handler.
var _ api.SecurityHandler = (*Handler)(nil)

// NewAuthHandler returns a new instance of the auth service handler.
func NewAuthHandler(auth *util.AuthService) *Handler {
	return &Handler{
		auth: auth,
	}
}

// HandleCookieAuth handles cookie based authentication.
func (h *Handler) HandleCookieAuth(ctx context.Context, _operationName string, t api.CookieAuth) (context.Context, error) {
	// Decrypt and read session cookie
	_, err := h.auth.ReadSession(ctx, t.APIKey)
	// If session does not exist, return error
	if err != nil {
		return nil, model.ErrUnauthorised
	}

	return ctx, nil
}
