package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(ctx context.Context, event *model.PageViewHit) error {
	exec := `--sql
	INSERT INTO views (bid, hostname, pathname, is_unique_user, is_unique_page, referrer, country_code, language, ua_browser, ua_os, ua_device_type, utm_source, utm_medium, utm_campaign, date_created)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())`

	_, err := c.DB.ExecContext(ctx, exec, event.BID, event.Hostname, event.Pathname, event.IsUniqueUser, event.IsUniquePage, event.Referrer, event.CountryCode, event.Language, event.BrowserName, event.OS, event.DeviceType, event.UTMSource, event.UTMMedium, event.UTMCampaign)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	_, err := c.DB.ExecContext(ctx, `--sql
		UPDATE views SET bid = NULL, duration_ms = ? WHERE bid = ?`,
		event.DurationMs, event.BID)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
