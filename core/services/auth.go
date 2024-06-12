package services

import (
	"context"
	"errors"
	"net/http"
	"time"

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

func (h *Handler) PostAuthLogout(ctx context.Context, params api.PostAuthLogoutParams) (api.PostAuthLogoutRes, error) {
	h.auth.RevokeSession(ctx, params.MeSess)

	// Expire cookie.
	cookie := &http.Cookie{
		Name:     model.SessionCookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	}

	return &api.PostAuthLogoutNoContent{
		SetCookie: cookie.String(),
	}, nil
}
