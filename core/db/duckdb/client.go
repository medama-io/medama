package duckdb

import (
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
		return nil, err
	}

	// Enable ICU extension
	_, err = db.Exec(`--sql
		INSTALL icu;`)
	if err != nil {
		return nil, err
	}

	// Load ICU extension
	_, err = db.Exec(`--sql
		LOAD icu;`)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
