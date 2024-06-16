package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(ctx context.Context, event *model.PageViewHit) error {
	exec := `--sql
	INSERT INTO views (
		bid,
		hostname,
		pathname,
		is_unique_user,
		is_unique_page,
		referrer,
		country_code,
		language,
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
		:referrer,
		:country_code,
		:language,
		:ua_browser,
		:ua_os,
		:ua_device_type,
		:utm_source,
		:utm_medium,
		:utm_campaign,
		NOW()
	)`

	paramMap := map[string]interface{}{
		"bid":            event.BID,
		"hostname":       event.Hostname,
		"pathname":       event.Pathname,
		"is_unique_user": event.IsUniqueUser,
		"is_unique_page": event.IsUniquePage,
		"referrer":       event.Referrer,
		"country_code":   event.CountryCode,
		"language":       event.Language,
		"ua_browser":     event.BrowserName,
		"ua_os":          event.OS,
		"ua_device_type": event.DeviceType,
		"utm_source":     event.UTMSource,
		"utm_medium":     event.UTMMedium,
		"utm_campaign":   event.UTMCampaign,
	}

	_, err := c.DB.NamedExecContext(ctx, exec, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	exec := `--sql
	UPDATE views SET duration_ms = :duration_ms WHERE bid = :bid`

	paramMap := map[string]interface{}{
		"bid":         event.BID,
		"duration_ms": event.DurationMs,
	}

	_, err := c.DB.NamedExecContext(ctx, exec, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
