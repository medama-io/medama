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

	bootQueries := []string{
		// Enable and load the ICU extension.
		"INSTALL icu;",
		"LOAD icu;",
	}

	for _, query := range bootQueries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, errors.Wrap(err, "duckdb")
		}
	}

	return &Client{
		DB: db,
	}, nil
}

// Close closes the database connection and any prepared statements.
func (c *Client) Close() error {
	// Helper function to close a statement and wrap any error.
	closeStmt := func(stmt *sqlx.NamedStmt) error {
		if stmt != nil {
			if err := stmt.Close(); err != nil {
				return errors.Wrap(err, "duckdb")
			}
		}
		return nil
	}

	// Close the statements.
	if err := closeStmt(addStmt); err != nil {
		return err
	}
	if err := closeStmt(updateStmt); err != nil {
		return err
	}

	// Close the database connection.
	return c.DB.Close()
}
