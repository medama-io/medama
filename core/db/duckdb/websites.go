package duckdb

import (
	"context"

	"github.com/go-faster/errors"
)

// DeleteWebsite deletes all rows associated with the given hostname.
func (c *Client) DeleteWebsite(ctx context.Context, hostname string) error {
	query := `--sql
		DELETE FROM views WHERE hostname = ?;`

	_, err := c.ExecContext(ctx, query, hostname)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
