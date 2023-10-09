package services

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) PostAuthLogin(ctx context.Context, req api.OptPostAuthLoginReq, params api.PostAuthLoginParams) (api.PostAuthLoginRes, error) {
	// Check session cookie and refresh if exists.
	sessionId := params.MeSess.Value
	if sessionId != "" {
		// Check if session exists.
		userId, err := h.auth.Cache.Get(ctx, sessionId)
		if err == nil {
			// Session exists, refresh token.
			cookie, err := h.auth.CreateSession(ctx, userId.(string), model.SessionDuration)
			if err != nil {
				return nil, err
			}

			return &api.PostAuthLoginOK{
				SetCookie: api.NewOptString(cookie.String()),
			}, nil
		}
	}

	// Check email and password.
	user, err := h.db.GetUserByEmail(ctx, req.Value.Email)
	if err != nil {
		return nil, err
	}

	match, err := h.auth.ComparePasswords(req.Value.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return ErrUnauthorised(err), nil
	}

	// Create session.
	cookie, err := h.auth.CreateSession(ctx, user.ID, model.SessionDuration)
	if err != nil {
		return nil, err
	}

	return &api.PostAuthLoginOK{
		SetCookie: api.NewOptString(cookie.String()),
	}, nil
}

func (h *Handler) PostAuthRefresh(ctx context.Context, req api.OptPostAuthRefreshReq, params api.PostAuthRefreshParams) (api.PostAuthRefreshRes, error) {
	// Check session cookie.
	sessionId := params.MeSess.Value
	if sessionId == "" {
		return ErrUnauthorised(nil), nil
	}

	// Check if session exists.
	userId, err := h.auth.Cache.Get(ctx, sessionId)
	if err != nil {
		return ErrUnauthorised(err), nil
	}

	// Refresh token.
	cookie, err := h.auth.CreateSession(ctx, userId.(string), model.SessionDuration)
	if err != nil {
		return nil, err
	}

	return &api.PostAuthRefreshOK{
		SetCookie: api.NewOptString(cookie.String()),
	}, nil
}
