package migrations

import (
	"github.com/medama-io/medama/db/sqlite"
)

func Up0006(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Update users table to create new settings JSON column.
	_, err = tx.Exec(`--sql
	ALTER TABLE users ADD COLUMN settings JSON NOT NULL DEFAULT '{}'`)
	if err != nil {
		return err
	}

	// Move language column to settings JSON column.
	_, err = tx.Exec(`--sql
	UPDATE users SET settings = JSON_OBJECT('language', language)`)
	if err != nil {
		return err
	}

	// Remove language column from users table.
	_, err = tx.Exec(`--sql
	ALTER TABLE users DROP COLUMN language`)
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

func Down0006(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Add language column back to users table.
	_, err = tx.Exec(`--sql
	ALTER TABLE users ADD COLUMN language TEXT NOT NULL`)
	if err != nil {
		return err
	}

	// Move language from settings JSON column back to language column.
	_, err = tx.Exec(`--sql
	UPDATE users SET language = JSON_EXTRACT(settings, '$.language')`)
	if err != nil {
		return err
	}

	// Remove settings JSON column from users table.
	_, err = tx.Exec(`--sql
	ALTER TABLE users DROP COLUMN settings`)
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
