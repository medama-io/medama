package sqlite

import (
	"fmt"

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
	// Enable foreign key support in sqlite
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s", host))
	if err != nil {
		return nil, errors.Wrap(err, "sqlite")
	}

	return &Client{
		DB: db,
	}, nil
}
