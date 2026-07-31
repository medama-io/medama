package migrations

import (
	"github.com/medama-io/medama/db/sqlite"
)

func Up0008(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Update users table to create new api key column.
	_, err = tx.Exec(`--sql
	ALTER TABLE users ADD COLUMN api_key TEXT default NULL`)
	if err != nil {
		return err
	}

	// Create unique index for api key column.
	_, err = tx.Exec(`--sql
	CREATE UNIQUE INDEX ux_users_api_key ON users(api_key)`)
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

func Down0008(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Remove unique index for api key column.
	_, err = tx.Exec(`--sql
	DROP INDEX IF EXISTS ux_users_api_key`)
	if err != nil {
		return err
	}

	// Remove api key column from users table.
	_, err = tx.Exec(`--sql
	ALTER TABLE users DROP COLUMN api_key`)
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
