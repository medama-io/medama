package services

import (
	"context"
	"errors"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) PostAuthLogin(ctx context.Context, req *api.AuthLogin) (api.PostAuthLoginRes, error) {
	// Check email and password.
	user, err := h.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return ErrNotFound(err), nil
		}

		return ErrInternalServerError(err), nil
	}

	// Compare password hashes.
	match, err := h.auth.ComparePasswords(req.Password, user.Password)
	if err != nil {
		return ErrInternalServerError(err), nil
	}
	if !match {
		return ErrUnauthorised(err), nil
	}

	// Create session.
	cookie, err := h.auth.CreateSession(ctx, user.ID)
	if err != nil {
		return ErrInternalServerError(err), nil
	}

	return &api.PostAuthLoginOK{
		SetCookie: cookie.String(),
	}, nil
}
