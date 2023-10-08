package services

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/util"
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

func (h *Handler) PostAuthLogin(ctx context.Context, req api.OptPostAuthLoginReq, params api.PostAuthLoginParams) (api.PostAuthLoginRes, error) {
	return nil, nil
}

func (h *Handler) PostAuthRefresh(ctx context.Context, req api.OptPostAuthRefreshReq, params api.PostAuthRefreshParams) (api.PostAuthRefreshRes, error) {
	return nil, nil
}
