package middlewares

import (
	"context"
	"errors"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"github.com/medama-io/medama/util/logger"
)

type Handler struct {
	auth    *util.AuthService
	apiKeys *util.ApiKeysService
}

// Compile time check for Handler.
var _ api.SecurityHandler = (*Handler)(nil)

// NewAuthHandler returns a new instance of the auth service handler.
func NewAuthHandler(
	auth *util.AuthService,
	apiKeys *util.ApiKeysService,
) *Handler {
	return &Handler{
		auth:    auth,
		apiKeys: apiKeys,
	}
}

// HandleCookieAuth handles cookie based authentication.
func (h *Handler) HandleCookieAuth(
	ctx context.Context,
	_operationName string,
	t api.CookieAuth,
) (context.Context, error) {
	// Decrypt and read session cookie
	userID, err := h.auth.ReadSession(ctx, t.APIKey)
	// If session does not exist, return error
	if err != nil {
		return nil, model.ErrUnauthorised
	}

	// We want to pass the validated user ID to the next handler
	ctx = context.WithValue(ctx, model.ContextKeyUserID, userID)

	return ctx, nil
}

func (h *Handler) HandleApiKey(
	ctx context.Context,
	_operationName string,
	t api.ApiKey,
) (context.Context, error) {
	user, err := h.apiKeys.ExchangeApiKey(ctx, t.APIKey)
	if err != nil {
		if errors.Is(err, model.ErrInvalidApiKey) {
			return nil, model.ErrUnauthorised
		}

		log := logger.Get()
		log.Err(err).Msg("failed to authorize using api key")

		return nil, model.ErrInternalServerError
	}

	ctx = context.WithValue(ctx, model.ContextKeyUserID, user.ID)

	return ctx, nil
}
