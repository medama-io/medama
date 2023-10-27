package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/medama-io/medama/model"
)

func (c *Client) CreateUser(ctx context.Context, user *model.User) error {
	exec := `--sql
	INSERT INTO users (id, email, password, language, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := c.DB.ExecContext(ctx, exec, user.ID, user.Email, user.Password, user.Language, user.DateCreated, user.DateUpdated)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if errors.Is(sqliteError.Code, sqlite3.ErrConstraint) {
				return model.ErrUserExists
			}
		}

		return err
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*model.User, error) {
	query := `--sql
	SELECT id, email, password, language, date_created, date_updated FROM users WHERE id = ?`

	res, err := c.DB.QueryxContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	defer res.Close()

	if res.Next() {
		user := &model.User{}

		err := res.StructScan(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, model.ErrUserNotFound
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `--sql
	SELECT id, email, password, language, date_created, date_updated FROM users WHERE email = ?`

	res, err := c.DB.QueryxContext(ctx, query, email)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next() {
		user := &model.User{}

		err := res.StructScan(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, model.ErrUserNotFound
}

func (c *Client) GetUserCount(ctx context.Context) (int64, error) {
	query := `--sql
	SELECT COUNT(*) FROM users`

	var count int64
	err := c.DB.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Client) UpdateUserEmail(ctx context.Context, id string, email string, dateUpdated int64) error {
	exec := `--sql
	UPDATE users SET email = ?, date_updated = ? WHERE id = ?`

	_, err := c.DB.ExecContext(ctx, exec, email, dateUpdated, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrUserNotFound
		}

		return err
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

		return err
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	exec := `--sql
	DELETE FROM users WHERE id = ?`

	res, err := c.DB.ExecContext(ctx, exec, id)
	if err != nil {
		return err
	}

	// Delete statement will silently succeed if the user does not exist.
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return model.ErrUserNotFound
	}

	return nil
}
