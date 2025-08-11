package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

const (
	addEventName       = "addEvent"
	addPageViewName    = "addPageView"
	updatePageViewName = "updatePageView"
)

// AddEvent adds an event with a custom property to the database.
func (c *Client) AddEvents(ctx context.Context, events *[]model.EventHit) error {
	return c.executeInTransaction(ctx, func(tx *sqlx.Tx) error {
		return c.addEventsWithinTransaction(ctx, tx, events)
	})
}

const addPageViewStmt = `--sql
		INSERT INTO views (
			bid,
			hostname,
			pathname,
			is_unique_user,
			is_unique_page,
			referrer_host,
			referrer_group,
			country,
			language_base,
			language_dialect,
			ua_browser,
			ua_os,
			ua_device_type,
			utm_source,
			utm_medium,
			utm_campaign,
			date_created
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			NOW()
		)`

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(
	ctx context.Context,
	event *model.PageViewHit,
	events *[]model.EventHit,
) error {
	return c.executeInTransaction(ctx, func(tx *sqlx.Tx) error {
		stmt, err := c.GetPreparedStmt(ctx, addPageViewName, addPageViewStmt)
		if err != nil {
			return errors.Wrap(err, "duckdb")
		}

		txStmt := tx.StmtxContext(ctx, stmt)

		_, err = txStmt.ExecContext(ctx,
			event.BID,
			event.Hostname,
			event.Pathname,
			event.IsUniqueUser,
			event.IsUniquePage,
			event.ReferrerHost,
			event.ReferrerGroup,
			event.Country,
			event.LanguageBase,
			event.LanguageDialect,
			event.BrowserName,
			event.OS,
			event.DeviceType,
			event.UTMSource,
			event.UTMMedium,
			event.UTMCampaign)
		if err != nil {
			return errors.Wrap(err, "duckdb: execute statement")
		}

		return c.addEventsWithinTransaction(ctx, tx, events)
	})
}

const updatePageViewStmt = `--sql
		UPDATE views SET duration_ms = ? WHERE bid = ?`

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	return c.executeInTransaction(ctx, func(tx *sqlx.Tx) error {
		stmt, err := c.GetPreparedStmt(ctx, updatePageViewName, updatePageViewStmt)
		if err != nil {
			return errors.Wrap(err, "duckdb")
		}

		txStmt := tx.StmtxContext(ctx, stmt)

		if _, err := txStmt.ExecContext(ctx, event.DurationMs, event.BID); err != nil {
			return errors.Wrap(err, "duckdb: execute statement")
		}

		return nil
	})
}

// executeInTransaction executes the given function within a transaction.
func (c *Client) executeInTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := c.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "duckdb: begin transaction")
	}
	defer tx.Rollback() //nolint: errcheck // Called on defer

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "duckdb: commit transaction")
	}

	return nil
}

const addEventStmt = `--sql
		INSERT INTO events (
			bid,
			batch_id,
			group_name,
			name,
			value,
			date_created
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			NOW()
		)`

// addEventsWithinTransaction adds events within an existing transaction.
func (c *Client) addEventsWithinTransaction(
	ctx context.Context,
	tx *sqlx.Tx,
	events *[]model.EventHit,
) error {
	if events == nil || len(*events) == 0 {
		return nil
	}

	stmt, err := c.GetPreparedStmt(ctx, addEventName, addEventStmt)
	if err != nil {
		return errors.Wrap(err, "duckdb: prepare statement")
	}

	txStmt := tx.StmtxContext(ctx, stmt)

	for _, event := range *events {
		_, err := txStmt.ExecContext(
			ctx,
			event.BID,
			event.BatchID,
			event.Group,
			event.Name,
			event.Value,
		)
		if err != nil {
			return errors.Wrap(err, "duckdb: execute statement")
		}
	}

	return nil
}
