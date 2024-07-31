package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

const (
	addEventName       = "addEvent"
	addPageViewName    = "addPageView"
	updatePageViewName = "updatePageView"
)

// AddEvent adds an event with a custom property to the database.
func (c *Client) AddEvents(ctx context.Context, events *[]model.EventHit) error {
	exec := `--sql
		INSERT INTO events (
			batch_id,
			group_name,
			name,
			value,
			date_created
		) VALUES (
			:batch_id,
			:group_name,
			:name,
			:value,
			NOW()
		)`

	stmt, err := c.GetPreparedStmt(ctx, addEventName, exec)
	if err != nil {
		return errors.Wrap(err, "duckdb")
	}

	// Start a transaction for batch insert
	tx, err := c.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "duckdb: begin event hit transaction")
	}
	defer tx.Rollback() //nolint: errcheck // Called on defer

	txStmt := tx.NamedStmtContext(ctx, stmt)

	for _, event := range *events {
		paramMap := map[string]interface{}{
			"batch_id":   event.BatchID,
			"group_name": event.Group,
			"name":       event.Name,
			"value":      event.Value,
		}

		_, err = txStmt.ExecContext(ctx, paramMap)
		if err != nil {
			return errors.Wrap(err, "duckdb")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "duckdb: commit event hit transaction")
	}

	return nil
}

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(ctx context.Context, event *model.PageViewHit) error {
	exec := `--sql
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
			:bid,
			:hostname,
			:pathname,
			:is_unique_user,
			:is_unique_page,
			:referrer_host,
			:referrer_group,
			:country,
			:language_base,
			:language_dialect,
			:ua_browser,
			:ua_os,
			:ua_device_type,
			:utm_source,
			:utm_medium,
			:utm_campaign,
			NOW()
		)`

	stmt, err := c.GetPreparedStmt(ctx, addPageViewName, exec)
	if err != nil {
		return errors.Wrap(err, "duckdb")
	}

	paramMap := map[string]interface{}{
		"bid":              event.BID,
		"hostname":         event.Hostname,
		"pathname":         event.Pathname,
		"is_unique_user":   event.IsUniqueUser,
		"is_unique_page":   event.IsUniquePage,
		"referrer_host":    event.ReferrerHost,
		"referrer_group":   event.ReferrerGroup,
		"country":          event.Country,
		"language_base":    event.LanguageBase,
		"language_dialect": event.LanguageDialect,
		"ua_browser":       event.BrowserName,
		"ua_os":            event.OS,
		"ua_device_type":   event.DeviceType,
		"utm_source":       event.UTMSource,
		"utm_medium":       event.UTMMedium,
		"utm_campaign":     event.UTMCampaign,
	}

	_, err = stmt.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	exec := `--sql
		UPDATE views SET duration_ms = :duration_ms WHERE bid = :bid`

	stmt, err := c.GetPreparedStmt(ctx, updatePageViewName, exec)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	paramMap := map[string]interface{}{
		"bid":         event.BID,
		"duration_ms": event.DurationMs,
	}

	_, err = stmt.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
