package sqlite

import (
	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/db"
)

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ db.AppClient = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	db, err := sqlx.Connect("sqlite3", host)
	if err != nil {
		return nil, errors.Wrap(err, "sqlite")
	}

	return &Client{
		DB: db,
	}, nil
}

// Close closes the database connection and any prepared statements.
func (c *Client) Close() error {
	return c.DB.Close()
}
