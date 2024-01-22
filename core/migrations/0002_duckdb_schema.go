package migrations

import (
	"fmt"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/model"
)

func Up0002(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// For reference, we sync these enums with the go-useragent library
	// github.com/medama-io/go-useragent

	// Create browser enum
	_, err = tx.Exec(fmt.Sprintf(`--sql
	CREATE TYPE browser AS ENUM ('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')`,
		model.UnknownBrowser,
		model.ChromeBrowser,
		model.EdgeBrowser,
		model.FirefoxBrowser,
		model.InternetExplorerBrowser,
		model.OperaBrowser,
		model.OperaMiniBrowser,
		model.SafariBrowser,
		model.VivaldiBrowser,
		model.SamsungBrowser,
		model.NintendoBrowser))

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	// Create os enum
	_, err = tx.Exec(fmt.Sprintf(`--sql
	CREATE TYPE os AS ENUM ('%s','%s','%s','%s','%s','%s','%s')`,
		model.UnknownOS,
		model.AndroidOS,
		model.ChromeOS,
		model.IOS,
		model.LinuxOS,
		model.MacOS,
		model.WindowsOS))

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	// Create device type enum
	_, err = tx.Exec(fmt.Sprintf(`--sql
	CREATE TYPE device_type AS ENUM ('%s','%s','%s','%s','%s')`,
		model.UnknownDevice,
		model.DesktopDevice,
		model.MobileDevice,
		model.TabletDevice,
		model.TVDevice))

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	// Create views table
	//
	// bid is beacon id used to link a page view and duration request together
	//
	// hostname is the hostname of the page view
	//
	// pathname is the pathname of the page view
	//
	// is_unique is true if the page view is unique
	//
	// referrer_hostname is the hostname of the referer http header
	//
	// referrer_pathname is the pathname of the referer http header
	//
	// country_code is the country code of the user
	//
	// language is the language code of the user
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
		pathname TEXT,
		is_unique BOOLEAN,
		referrer_hostname TEXT,
		referrer_pathname TEXT,
		country_code TEXT,
		language TEXT,
		ua_browser browser,
		ua_os os,
		ua_device_type device_type,
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
