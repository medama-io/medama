package services

import "github.com/medama-io/medama/db/sqlite"

type Handler struct {
	db *sqlite.Client
}

func NewService(db *sqlite.Client) *Handler {
	return &Handler{
		db: db,
	}
}
