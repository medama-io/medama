package sqlite

import (
	"fmt"

	"github.com/medama-io/medama/db"

	"github.com/jmoiron/sqlx"
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
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
