package migrations

import (
	"github.com/medama-io/medama/db/duckdb"
)

func Up0004(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Create events table
	//
	// group is the group name of the event, typically the hostname
	//
	// name is the name of the event
	//
	// value is the value of the event
	//
	// date_created is the date the event was created
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS events (
		group_name TEXT NOT NULL,
		name TEXT NOT NULL,
		value TEXT NOT NULL,
		date_created TIMESTAMPTZ NOT NULL
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

func Down0004(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Drop views table
	_, err = tx.Exec(`--sql
	DROP TABLE IF EXISTS events`)
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
