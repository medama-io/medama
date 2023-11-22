package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

func (c *Client) AddPageView(ctx context.Context, event *model.PageView) error {
	exec := `--sql
	INSERT INTO views (bid, hostname, pathname, is_unique, referrer, title, timezone, language, ua_raw, ua_browser, ua_version, ua_os, ua_device_type, screen_width, screen_height, date_created)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Insert the page view into the database
	_, err := c.DB.ExecContext(ctx, exec, event.BID, event.Hostname, event.Pathname, event.IsUnique, event.Referrer, event.Title, event.Timezone, event.Language, event.RawUserAgent, event.BrowserName, event.BrowserVersion, event.OS, event.DeviceType, event.ScreenWidth, event.ScreenHeight, event.DateCreated)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewUpdate) error {
	// Update the page view into the database
	_, err := c.DB.ExecContext(ctx, `--sql
		UPDATE views SET duration_ms = ? WHERE bid = ?`,
		event.DurationMs, event.BID)
	if err != nil {
		return err
	}

	return nil
}
