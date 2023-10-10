package services

import (
	"context"

	"github.com/medama-io/medama/api"
)

func (h *Handler) PostAuthLogin(ctx context.Context, req api.OptPostAuthLoginReq) (api.PostAuthLoginRes, error) {
	// Check email and password.
	user, err := h.db.GetUserByEmail(ctx, req.Value.Email)
	if err != nil {
		return nil, err
	}

	// Compare password hashes.
	match, err := h.auth.ComparePasswords(req.Value.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return ErrUnauthorised(err), nil
	}

	// Create session.
	cookie, err := h.auth.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &api.PostAuthLoginOK{
		SetCookie: api.NewOptString(cookie.String()),
	}, nil
}
