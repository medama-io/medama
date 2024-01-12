package migrations

import (
	"context"
	"time"

	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"go.jetpack.io/typeid"
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
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		language TEXT NOT NULL,
		date_created INTEGER NOT NULL,
		date_updated INTEGER NOT NULL,
		UNIQUE(username)
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

	// Create default admin user
	// UUIDv7 id generation
	typeid, err := typeid.WithPrefix("user")
	if err != nil {
		return err
	}
	id := typeid.String()

	// Hash default password
	auth, err := util.NewAuthService(context.Background())
	if err != nil {
		return err
	}
	pwdHash, err := auth.HashPassword("admin")
	if err != nil {
		return err
	}

	dateCreated := time.Now().Unix()
	dateUpdated := dateCreated
	err = c.CreateUser(context.Background(), model.NewUser(id, "admin", pwdHash, "en", dateCreated, dateUpdated))
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
	_, err = tx.Exec(`--sql
	DROP TABLE IF EXISTS users`)
	if err != nil {
		return err
	}

	// Drop websites table
	_, err = tx.Exec(`--sql
	DROP TABLE IF EXISTS websites`)
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
