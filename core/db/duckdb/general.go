package duckdb

import (
	"context"
)

func (c *Client) GetDatabaseVersion(ctx context.Context) (string, error) {
	_, err := c.DB.ExecContext(ctx, "PRAGMA version;")
	if err != nil {
		return "", err
	}

	var version struct {
		Version  string `db:"library_version"`
		SourceID string `db:"source_id"`
	}

	err = c.DB.GetContext(ctx, &version, "CALL pragma_version();")
	if err != nil {
		return "", err
	}

	return version.Version, nil
}
