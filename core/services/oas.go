package services

import (
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/util"
)

type Handler struct {
	auth *AuthService
	db   *sqlite.Client
}

// NewService returns a new instance of the ogen service handler.
func NewService(cache *util.Cache, db *sqlite.Client) *Handler {
	return &Handler{
		auth: NewAuthService(cache),
		db:   db,
	}
}
