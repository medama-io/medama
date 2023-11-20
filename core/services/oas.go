package services

import (
	"github.com/medama-io/go-useragent"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/util"
)

type Handler struct {
	auth        *util.AuthService
	db          *sqlite.Client
	analyticsDB *duckdb.Client
	useragent   *useragent.Parser
}

// NewService returns a new instance of the ogen service handler.
func NewService(auth *util.AuthService, sqlite *sqlite.Client, duckdb *duckdb.Client) *Handler {
	return &Handler{
		auth:        auth,
		db:          sqlite,
		analyticsDB: duckdb,
		useragent:   useragent.NewParser(),
	}
}
