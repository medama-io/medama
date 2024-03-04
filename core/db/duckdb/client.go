package duckdb

import (
	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/db"
)

type Client struct {
	*sqlx.DB
}

// Compile time check for Client.
var (
	_ db.AnalyticsClient = (*Client)(nil)
)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	db, err := sqlx.Connect("duckdb", host)
	if err != nil {
		return nil, errors.Wrap(err, "duckdb")
	}

	// Enable ICU extension
	_, err = db.Exec(`--sql
		INSTALL icu;`)
	if err != nil {
		return nil, errors.Wrap(err, "duckdb")
	}

	// Load ICU extension
	_, err = db.Exec(`--sql
		LOAD icu;`)
	if err != nil {
		return nil, errors.Wrap(err, "duckdb")
	}

	return &Client{
		DB: db,
	}, nil
}
