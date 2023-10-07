package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/medama-io/medama/model"
)

func (c *Client) CreateUser(ctx context.Context, user *model.User) error {
	_, err := c.DB.ExecContext(ctx, "INSERT INTO users (id, email, password, language, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Email, user.Password, user.Language, user.DateCreated, user.DateUpdated)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*model.User, error) {
	res, err := c.DB.QueryxContext(ctx, "SELECT id, email, password, language, date_created, date_updated FROM users WHERE id = ?", id)
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

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	res, err := c.DB.QueryxContext(ctx, "SELECT id, email, password, language, date_created, date_updated FROM users WHERE email = ?", email)
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

func (c *Client) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := c.DB.ExecContext(ctx, "UPDATE users SET email = ?, password = ?, language = ?, date_created = ?, date_updated = ? WHERE id = ?", user.Email, user.Password, user.Language, user.DateCreated, user.DateUpdated, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	_, err := c.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrUserNotFound
		}

		return err
	}

	return nil
}
