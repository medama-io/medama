package migrations

import (
	"github.com/medama-io/medama/db/duckdb"
)

func Up0005(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Update events table to include optional bid column.
	_, err = tx.Exec(`--sql
	ALTER TABLE events ADD COLUMN bid TEXT`)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Down0005(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Drop bid column from events table.
	_, err = tx.Exec(`--sql
	ALTER TABLE events DROP bid`)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
