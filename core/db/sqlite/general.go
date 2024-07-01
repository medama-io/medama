package sqlite

import (
	"context"
)

func (c *Client) GetDatabaseVersion(ctx context.Context) (string, error) {
	query := `--sql
		select sqlite_version();`

	var version string
	err := c.DB.GetContext(ctx, &version, query)
	if err != nil {
		return "", err
	}

	return "v" + version, nil
}
