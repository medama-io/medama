package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
	"github.com/ncruces/go-sqlite3"
)

func (c *Client) CreateUser(ctx context.Context, user *model.User) error {
	exec := `--sql
	INSERT INTO users (
		id,
		username,
		password,
		settings,
		date_created,
		date_updated
	) VALUES (
		:id,
		:username,
		:password,
		:settings,
		:date_created,
		:date_updated
	)`

	// Marshal settings to JSON
	settingsJSON, err := json.Marshal(user.Settings)
	if err != nil {
		return errors.Wrap(err, "marshaling settings")
	}

	paramMap := map[string]any{
		"id":           user.ID,
		"username":     user.Username,
		"password":     user.Password,
		"settings":     string(settingsJSON),
		"date_created": user.DateCreated,
		"date_updated": user.DateUpdated,
	}

	_, err = c.DB.NamedExecContext(ctx, exec, paramMap)
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
	SELECT id, username, password, settings, date_created, date_updated FROM users WHERE id = ?`

	var user model.User
	var settingsJSON string

	err := c.DB.QueryRowxContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&settingsJSON,
		&user.DateCreated,
		&user.DateUpdated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "db")
	}

	// Parse the JSON settings
	if settingsJSON != "" {
		user.Settings = model.NewDefaultSettings()
		err = json.Unmarshal([]byte(settingsJSON), user.Settings)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal settings")
		}
	}

	return &user, nil
}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `--sql
	SELECT id, username, password, settings, date_created, date_updated FROM users WHERE username = ?`

	var user model.User
	var settingsJSON string

	err := c.DB.QueryRowxContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&settingsJSON,
		&user.DateCreated,
		&user.DateUpdated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "db")
	}

	// Parse the JSON settings
	if settingsJSON != "" {
		user.Settings = model.NewDefaultSettings()
		err = json.Unmarshal([]byte(settingsJSON), user.Settings)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal settings")
		}
	}

	return &user, nil
}

func (c *Client) UpdateUserUsername(ctx context.Context, id string, username string) error {
	exec := `--sql
	UPDATE users SET username = :username, date_updated = :date_updated WHERE id = :id`

	paramMap := map[string]any{
		"id":           id,
		"username":     username,
		"date_updated": time.Now().Unix(),
	}

	_, err := c.DB.NamedExecContext(ctx, exec, paramMap)
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

func (c *Client) UpdateUserPassword(ctx context.Context, id string, password string) error {
	exec := `--sql
	UPDATE users SET password = :password, date_updated = :date_updated WHERE id = :id`

	paramMap := map[string]any{
		"id":           id,
		"password":     password,
		"date_updated": time.Now().Unix(),
	}

	_, err := c.DB.NamedExecContext(ctx, exec, paramMap)
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
