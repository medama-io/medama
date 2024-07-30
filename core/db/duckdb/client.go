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
	statements *haxmap.Map[string, *sqlx.NamedStmt]
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

	prepared := haxmap.New[string, *sqlx.NamedStmt]()
	c := &Client{
		DB:         db,
		statements: prepared,
	}

	err = c.prepareStatements(context.Background())
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Close closes the database connection and any prepared statements.
func (c *Client) Close() error {
	// Close the statements.
	c.closeStatements(context.Background())

	// Close the database connection.
	return c.DB.Close()
}

func (c *Client) prepareStatements(ctx context.Context) error {
	// Map of statements to prepare on startup
	queryMap := map[string]string{
		addEventName:       addEventQuery,
		addPageViewName:    addPageViewQuery,
		updatePageViewName: updatePageViewQuery,
	}

	for name := range queryMap {
		stmt, err := c.DB.PrepareNamedContext(ctx, queryMap[name])
		if err != nil {
			return errors.Wrap(err, "duckdb: unable to create prepared statement")
		}
		c.statements.Set(name, stmt)
	}

	return nil
}

func (c *Client) closeStatements(ctx context.Context) {
	c.statements.ForEach(func(_ string, stmt *sqlx.NamedStmt) bool {
		stmt.Close()
		return true
	})
}
