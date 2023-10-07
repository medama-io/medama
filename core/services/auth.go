package services

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/medama-io/medama/api"
)

type AuthService struct{}

// NewAuthService returns a new instance of AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
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
