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
	INSERT INTO websites (id, user_id, hostname, date_created, date_updated) VALUES (?, ?, ?, ?, ?)`

	_, err := c.DB.ExecContext(ctx, exec, website.ID, website.UserID, website.Hostname, website.DateCreated, website.DateUpdated)
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
			slog.String("id", website.ID),
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
	SELECT id, user_id, hostname, date_created, date_updated FROM websites WHERE user_id = ?`

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

func (c *Client) GetWebsite(ctx context.Context, id string) (*model.Website, error) {
	var website model.Website

	query := `--sql
	SELECT id, user_id, hostname, date_created, date_updated FROM websites WHERE id = ?`

	err := c.QueryRowxContext(ctx, query, id).StructScan(&website)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.DebugContext(ctx, "website not found", slog.String("id", id))
			return nil, model.ErrWebsiteNotFound
		}

		attributes := []slog.Attr{
			slog.String("id", id),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get website", attributes...)

		return nil, err
	}

	return &website, nil
}

func (c *Client) DeleteWebsite(ctx context.Context, id string) error {
	exec := `--sql
	DELETE FROM websites WHERE id = ?`

	res, err := c.DB.ExecContext(ctx, exec, id)
	if err != nil {
		attributes := []slog.Attr{
			slog.String("id", id),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to delete website", attributes...)

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		attributes := []slog.Attr{
			slog.String("id", id),
			slog.String("error", err.Error()),
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get rows affected", attributes...)

		return err
	}

	if rowsAffected == 0 {
		slog.DebugContext(ctx, "website not found", slog.String("id", id))
		return model.ErrWebsiteNotFound
	}

	return nil
}
