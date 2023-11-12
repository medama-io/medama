package migrations

import (
	"github.com/medama-io/medama/db/duckdb"
)

func Up0002(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Create views table
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS views (
		bid TEXT PRIMARY KEY,
		hostname TEXT NOT NULL,
		pathname TEXT,
		referrer TEXT,
		title TEXT,
		timezone TEXT,
		language TEXT,
		screen_width INTEGER,
		screen_height INTEGER,
		duration_ms INTEGER,
		date_created INTEGER NOT NULL
	)`)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Down0002(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Drop views table
	_, err = tx.Exec(`--sql
	DROP TABLE IF EXISTS views`)
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
