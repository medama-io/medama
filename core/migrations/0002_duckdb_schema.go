// nolint: wrapcheck
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
	//
	// bid is the unique beacon id of the page view that links to the unload beacon
	//
	// hostname is the hostname of the page view
	//
	// pathname is the pathname of the page view
	//
	// is_unique_user is true if the user is unique
	//
	// is_unique_page is true if the user has not visited the page before
	//
	// referrer_host is the referer hostname from the http header of the page view
	//
	// referrer_group is the referer group of the page view
	//
	// country is the country code of the user
	//
	// language_base is the base language of the user e.g. English
	//
	// language_dialect is the dialect of the user's language e.g. British English
	//
	// ua_browser is the browser of the user
	//
	// ua_os is the operating system of the user
	//
	// ua_device_type is the device type of the user
	//
	// utm_source is the utm source of the user
	//
	// utm_medium is the utm medium of the user
	//
	// utm_campaign is the utm campaign of the user
	//
	// duration_ms is the duration of the page view in milliseconds
	//
	// date_created is the date the page view was created
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS views (
		bid TEXT,
		hostname TEXT NOT NULL,
		pathname TEXT NOT NULL,
		is_unique_user BOOLEAN NOT NULL,
		is_unique_page BOOLEAN NOT NULL,
		referrer_host TEXT,
		referrer_group TEXT,
		country TEXT,
		language_base TEXT,
		language_dialect TEXT,
		ua_browser TEXT NOT NULL,
		ua_os TEXT NOT NULL,
		ua_device_type TEXT NOT NULL,
		utm_source TEXT,
		utm_medium TEXT,
		utm_campaign TEXT,
		duration_ms UINTEGER,
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
