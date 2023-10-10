package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/mattn/go-sqlite3"
	"github.com/medama-io/medama/model"
)

func (c *Client) CreateWebsite(ctx context.Context, website *model.Website) error {
	exec := `--sql
	INSERT INTO websites (user_id, hostname, name, date_created, date_updated) VALUES (?, ?, ?, ?, ?)`

	_, err := c.DB.ExecContext(ctx, exec, website.UserID, website.Hostname, website.Name, website.DateCreated, website.DateUpdated)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			switch sqliteError.ExtendedCode {
			// Check for unique hostname constraint
			case sqlite3.ErrConstraintPrimaryKey:
				return model.ErrWebsiteExists

				// Check for foreign key constraint
			case sqlite3.ErrConstraintForeignKey:
				return model.ErrUserNotFound
			}
		}

		attributes := []slog.Attr{
			slog.String("user_id", website.UserID),
			slog.String("hostname", website.Hostname),
			slog.Int64("date_created", website.DateCreated),
			slog.Int64("date_updated", website.DateUpdated),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to create website", attributes...)

		return err
	}

	return nil
}

func (c *Client) ListWebsites(ctx context.Context, userID string) ([]*model.Website, error) {
	var websites []*model.Website

	query := `--sql
	SELECT user_id, hostname, name, date_created, date_updated FROM websites WHERE user_id = ?`

	err := c.SelectContext(ctx, &websites, query, userID)
	if err != nil {
		attributes := []slog.Attr{
			slog.String("user_id", userID),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to list websites", attributes...)

		return nil, err
	}

	if len(websites) == 0 {
		slog.DebugContext(ctx, "no websites found", slog.String("user_id", userID))
		return nil, model.ErrWebsiteNotFound
	}

	return websites, nil
}

func (c *Client) UpdateWebsite(ctx context.Context, website *model.Website) error {
	// Update all values except user_id
	exec := `--sql
	UPDATE websites SET hostname = ?, name = ?, date_updated = ? WHERE hostname = ?`

	res, err := c.DB.ExecContext(ctx, exec, website.Hostname, website.Name, website.DateUpdated, website.Hostname)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			// Check for unique hostname constraint
			if sqliteError.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return model.ErrWebsiteExists
			}
		}

		attributes := []slog.Attr{
			slog.String("hostname", website.Hostname),
			slog.String("name", website.Name),
			slog.Int64("date_updated", website.DateUpdated),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to update website", attributes...)

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		attributes := []slog.Attr{
			slog.String("hostname", website.Hostname),
			slog.String("name", website.Name),
			slog.Int64("date_updated", website.DateUpdated),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get rows affected", attributes...)

		return err
	}

	if rowsAffected == 0 {
		slog.DebugContext(ctx, "website not found", slog.String("hostname", website.Hostname))
		return model.ErrWebsiteNotFound
	}

	return nil
}

func (c *Client) GetWebsite(ctx context.Context, hostname string) (*model.Website, error) {
	var website model.Website

	query := `--sql
	SELECT user_id, hostname, name, date_created, date_updated FROM websites WHERE hostname = ?`

	err := c.QueryRowxContext(ctx, query, hostname).StructScan(&website)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.DebugContext(ctx, "website not found", slog.String("hostname", hostname))
			return nil, model.ErrWebsiteNotFound
		}

		attributes := []slog.Attr{
			slog.String("hostname", hostname),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get website", attributes...)

		return nil, err
	}

	return &website, nil
}

func (c *Client) DeleteWebsite(ctx context.Context, hostname string) error {
	exec := `--sql
	DELETE FROM websites WHERE hostname = ?`

	res, err := c.DB.ExecContext(ctx, exec, hostname)
	if err != nil {
		attributes := []slog.Attr{
			slog.String("hostname", hostname),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to delete website", attributes...)

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		attributes := []slog.Attr{
			slog.String("hostname", hostname),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get rows affected", attributes...)

		return err
	}

	if rowsAffected == 0 {
		slog.DebugContext(ctx, "website not found", slog.String("hostname", hostname))
		return model.ErrWebsiteNotFound
	}

	return nil
}
