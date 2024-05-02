package sqlite

import (
	"context"
	"database/sql"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
	"github.com/ncruces/go-sqlite3"
)

func (c *Client) CreateUser(ctx context.Context, user *model.User) error {
	exec := `--sql
	INSERT INTO users (id, username, password, language, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := c.DB.ExecContext(ctx, exec, user.ID, user.Username, user.Password, user.Language, user.DateCreated, user.DateUpdated)
	if err != nil {
		if errors.Is(err, sqlite3.CONSTRAINT_UNIQUE) || errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY) {
			return model.ErrUserExists
		}

		return errors.Wrap(err, "db")
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*model.User, error) {
	query := `--sql
	SELECT id, username, password, language, date_created, date_updated FROM users WHERE id = ?`

	res, err := c.DB.QueryxContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "db")
	}

	defer res.Close()

	if res.Next() {
		user := &model.User{}

		err := res.StructScan(user)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}

		return user, nil
	}

	return nil, model.ErrUserNotFound
}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `--sql
	SELECT id, username, password, language, date_created, date_updated FROM users WHERE username = ?`

	res, err := c.DB.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	defer res.Close()

	if res.Next() {
		user := &model.User{}

		err := res.StructScan(user)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}

		return user, nil
	}

	return nil, model.ErrUserNotFound
}

func (c *Client) UpdateUserUsername(ctx context.Context, id string, username string, dateUpdated int64) error {
	exec := `--sql
	UPDATE users SET username = ?, date_updated = ? WHERE id = ?`

	_, err := c.DB.ExecContext(ctx, exec, username, dateUpdated, id)
	if err != nil {
		switch {
		case errors.Is(err, sqlite3.CONSTRAINT_UNIQUE),
			errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY):
			return model.ErrUserExists
		case errors.Is(err, sql.ErrNoRows):
			return model.ErrUserNotFound
		}
		return errors.Wrap(err, "db")
	}

	return nil
}

func (c *Client) UpdateUserPassword(ctx context.Context, id string, password string, dateUpdated int64) error {
	exec := `--sql
	UPDATE users SET password = ?, date_updated = ? WHERE id = ?`

	_, err := c.DB.ExecContext(ctx, exec, password, dateUpdated, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrUserNotFound
		}

		return errors.Wrap(err, "db")
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	exec := `--sql
	DELETE FROM users WHERE id = ?`

	res, err := c.DB.ExecContext(ctx, exec, id)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	// Delete statement will silently succeed if the user does not exist.
	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "db")
	}

	if count == 0 {
		return model.ErrUserNotFound
	}

	return nil
}
