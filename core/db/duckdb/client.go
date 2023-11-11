package duckdb

import (
	"github.com/jmoiron/sqlx"
)

type Handler interface {
	// Events
}

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ Handler = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	// Enable foreign key support in sqlite
	db, err := sqlx.Connect("duckdb", host)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
