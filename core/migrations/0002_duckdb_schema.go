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
		is_unique BOOLEAN,
		referrer_hostname TEXT,
		referrer_pathname TEXT,
		title TEXT,
		country_code TEXT,
		language TEXT,
		ua_raw TEXT,
		ua_browser UTINYINT,
		ua_version TEXT,
		ua_os UTINYINT,
		ua_device_type UTINYINT,
		screen_width USMALLINT,
		screen_height USMALLINT,
		utm_source TEXT,
		utm_medium TEXT,
		utm_campaign TEXT,
		duration_ms UINTEGER,
		date_updated TIMESTAMPTZ NOT NULL
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
