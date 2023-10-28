package migrations

import (
	"github.com/medama-io/medama/db/sqlite"
)

func Up0001(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Create users table
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		language TEXT NOT NULL,
		date_created INTEGER NOT NULL,
		date_updated INTEGER NOT NULL,
		UNIQUE(email)
	)`)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	// Create websites table
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS websites (
		hostname TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		name TEXT NOT NULL,
		date_created INTEGER NOT NULL,
		date_updated INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
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

func Down0001(c *sqlite.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Drop users table
	_, err = tx.Exec("DROP TABLE IF EXISTS users")
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
