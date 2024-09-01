package duckdb

import (
	"context"

	"github.com/alphadose/haxmap"
	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/db"
)

type Client struct {
	*sqlx.DB
	// Map of prepared statements.
	statements *haxmap.Map[string, *sqlx.Stmt]
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
		DB:         db,
		statements: haxmap.New[string, *sqlx.Stmt](),
	}, nil
}

// Close closes the database connection and any prepared statements.
func (c *Client) Close() error {
	// Close the statements.
	c.closeStatements()

	// Close the database connection.
	return c.DB.Close()
}

// GetPreparedStmt returns a prepared statement by name. This is lazy loaded and cached after
// the first call.
func (c *Client) GetPreparedStmt(ctx context.Context, name string, query string) (*sqlx.Stmt, error) {
	stmt, ok := c.statements.Get(name)
	if ok {
		return stmt, nil
	}

	stmt, err := c.DB.PreparexContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create prepared statement")
	}

	c.statements.Set(name, stmt)
	return stmt, nil
}

func (c *Client) closeStatements() {
	c.statements.ForEach(func(_ string, stmt *sqlx.Stmt) bool {
		stmt.Close()
		return true
	})
}
