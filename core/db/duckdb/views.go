package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

func (c *Client) AddPageView(ctx context.Context, event *model.PageView) error {
	// Insert the page view into the database
	_, err := c.DB.ExecContext(ctx, `--sql
		INSERT INTO views (bid, event_type, hostname, pathname, referrer, title, timezone, screen_width, screen_height, duration_ms, date_created)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.BID, event.EventType, event.Hostname, event.Pathname, event.Referrer, event.Title, event.Timezone, event.ScreenWidth, event.ScreenHeight, event.DurationMs, event.DateCreated)
	if err != nil {
		return err
	}

	return nil
}
