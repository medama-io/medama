package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

func (c *Client) AddPageView(ctx context.Context, event *model.PageView) error {
	// Insert the page view into the database
	_, err := c.DB.ExecContext(ctx, `--sql
		INSERT INTO views (bid, hostname, pathname, referrer, title, timezone, language, screen_width, screen_height, date_created)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.BID, event.Hostname, event.Pathname, event.Referrer, event.Title, event.Timezone, event.Language, event.ScreenWidth, event.ScreenHeight, event.DateCreated)
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
