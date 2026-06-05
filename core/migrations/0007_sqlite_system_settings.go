package migrations

import (
	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db/sqlite"
)

func Up0007(c *sqlite.Client) error {
	tx, err := c.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin migration")
	}

	// Create system_settings table
	_, err = tx.Exec(`--sql
	CREATE TABLE IF NOT EXISTS system_settings (
		key TEXT NOT NULL PRIMARY KEY,
		value TEXT NOT NULL,
		date_updated INTEGER NOT NULL
	)`)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(
				errors.Join(err, rollbackErr),
				"failed to create system_settings table",
			)
		}

		return errors.Wrap(err, "failed to create system_settings table")
	}

	// Migrate script_type, block_abusive_ips, block_tor_exit_nodes, blocked_ips
	// from users.settings (JSON) to system_settings table.
	// Keep only language in users.settings.
	for _, key := range []string{
		"script_type",
		"block_abusive_ips",
		"block_tor_exit_nodes",
		"blocked_ips",
	} {
		_, err = tx.Exec(`--sql
		INSERT OR REPLACE INTO system_settings (key, value, date_updated)
		SELECT ?, JSON_EXTRACT(settings, '$.' || ?), CAST(strftime('%s', 'now') AS INTEGER)
		FROM users
		WHERE JSON_EXTRACT(settings, '$.' || ?) IS NOT NULL`, key, key, key)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return errors.Wrap(
					errors.Join(err, rollbackErr),
					"failed to migrate settings to system_settings table",
				)
			}

			return errors.Wrap(err, "failed to migrate settings to system_settings table")
		}
	}

	// Remove migrated keys from users.settings, keeping only language
	_, err = tx.Exec(`--sql
	UPDATE users SET
		settings = JSON_REMOVE(
			settings,
			'$.script_type',
			'$.block_abusive_ips',
			'$.block_tor_exit_nodes',
			'$.blocked_ips'
		)`)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(
				errors.Join(err, rollbackErr),
				"failed to cleanup user settings",
			)
		}

		return errors.Wrap(err, "failed to cleanup user settings")
	}

	return tx.Commit()
}

func Down0007(c *sqlite.Client) error {
	tx, err := c.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin migration rollback")
	}

	// Restore system settings back into users.settings JSON
	_, err = tx.Exec(`--sql
	UPDATE users SET
		settings = JSON_SET(
			settings,
			'$.script_type', (SELECT value FROM system_settings WHERE key = 'script_type'),
			'$.block_abusive_ips', (SELECT value FROM system_settings WHERE key = 'block_abusive_ips'),
			'$.block_tor_exit_nodes', (SELECT value FROM system_settings WHERE key = 'block_tor_exit_nodes'),
			'$.blocked_ips', (SELECT value FROM system_settings WHERE key = 'blocked_ips')
		)`)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(
				errors.Join(err, rollbackErr),
				"failed to restore user settings",
			)
		}

		return errors.Wrap(err, "failed to restore user settings")
	}

	// Drop system_settings table
	_, err = tx.Exec(`--sql
	DROP TABLE IF EXISTS system_settings`)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(
				errors.Join(err, rollbackErr),
				"failed to remove system_settings table",
			)
		}

		return errors.Wrap(err, "failed to remove system_settings table")
	}

	return tx.Commit()
}
