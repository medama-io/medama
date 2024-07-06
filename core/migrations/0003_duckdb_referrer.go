package migrations

import (
	"github.com/medama-io/go-referrer-parser"
	"github.com/medama-io/medama/db/duckdb"
)

func Up0003(c *duckdb.Client) error {
	// Start referrer parser
	parser, err := referrer.NewParser()
	if err != nil {
		return err
	}

	// Get all distinct referrer hostnames
	rows, err := c.Queryx(`--sql
		SELECT DISTINCT referrer_host FROM views
	`)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Iterate over all referrer hostnames and parse them. If the referrer returns
	// not nil, update all rows with the same referrer hostnames with the parsed
	// referrer group.
	for rows.Next() {
		var referrerHost string
		err = rows.Scan(&referrerHost)
		if err != nil {
			return err
		}

		referrerGroup := parser.Parse(referrerHost)
		if referrerGroup != "" {
			_, err = tx.Exec(`--sql
				UPDATE views SET referrer_group = ? WHERE referrer_host = ?
			`, referrerGroup, referrerHost)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Down0003(c *duckdb.Client) error {
	// Begin transaction
	tx, err := c.Beginx()
	if err != nil {
		return err
	}

	// Set all referrer groups to empty string
	_, err = tx.Exec(`--sql
		UPDATE views SET referrer_group = ''
	`)
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
