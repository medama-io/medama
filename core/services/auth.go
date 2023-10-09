package services

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"go.jetpack.io/typeid"
)

type AuthService struct {
	// Cache used to store session tokens.
	cache *util.Cache
}

// NewAuthService returns a new instance of AuthService.
func NewAuthService(cache *util.Cache) *AuthService {
	return &AuthService{
		cache: cache,
	}
}

// HashPassword hashes a password using argon.
func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// ComparePasswords compares a supplied password with a stored hash.
func (a *AuthService) ComparePasswords(suppliedPassword string, storedHash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(suppliedPassword, storedHash)
	if err != nil {
		return false, err
	}

	return match, nil
}

// CreateSession creates a new session token and stores it in the cache.
func (a *AuthService) CreateSession(ctx context.Context, userId string, duration time.Duration) (string, *http.Cookie, error) {
	// Generate session token.
	sessionIdType, err := typeid.New("sess")
	if err != nil {
		return "", nil, err
	}
	sessionId := sessionIdType.String()

	// Set session token in cache.
	a.cache.Set(sessionId, userId, duration)

	// Create session cookie.
	cookie := &http.Cookie{
		Name:     model.SessionCookieName,
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return sessionId, cookie, nil
}

// RevokeSession deletes a session token from the cache.
func (a *AuthService) RevokeSession(ctx context.Context, sessionId string) {
	a.cache.Delete(sessionId)
}

func (h *Handler) PostAuthLogin(ctx context.Context, req api.OptPostAuthLoginReq, params api.PostAuthLoginParams) (api.PostAuthLoginRes, error) {
	// Check session cookie and refresh if exists.
	sessionId := params.MeSess.Value
	if sessionId != "" {
		// Check if session exists.
		userId, err := h.auth.cache.Get(ctx, sessionId)
		if err == nil {
			// Session exists, refresh token.
			_, cookie, err := h.auth.CreateSession(ctx, userId.(string), model.SessionDuration)
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
	_, cookie, err := h.auth.CreateSession(ctx, user.ID, model.SessionDuration)
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
	userId, err := h.auth.cache.Get(ctx, sessionId)
	if err != nil {
		return ErrUnauthorised(err), nil
	}

	// Refresh token.
	_, cookie, err := h.auth.CreateSession(ctx, userId.(string), model.SessionDuration)
	if err != nil {
		return nil, err
	}

	return &api.PostAuthRefreshOK{
		SetCookie: api.NewOptString(cookie.String()),
	}, nil
}
