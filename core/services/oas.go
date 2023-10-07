package services

import "github.com/medama-io/medama/db/sqlite"

type Handler struct {
	auth *AuthService
	db   *sqlite.Client
}

func NewService(db *sqlite.Client) *Handler {
	return &Handler{
		auth: NewAuthService(),
		db:   db,
	}
}
