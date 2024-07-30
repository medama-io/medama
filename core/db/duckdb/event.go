package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

const (
	addEventName  = "addEvent"
	addEventQuery = `--sql
		INSERT INTO events (
			group,
			name,
			value,
			date_created
		) VALUES (
			:group,
			:name,
			:value,
			NOW()
		)`
	addPageViewName  = "addPageView"
	addPageViewQuery = `--sql
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
	updatePageViewName  = "updatePageView"
	updatePageViewQuery = `--sql
		UPDATE views SET duration_ms = :duration_ms WHERE bid = :bid`
)

// AddEvent adds an event with a custom property to the database.
func (c *Client) AddEvent(ctx context.Context, event *model.EventHit) error {
	paramMap := map[string]interface{}{
		"group": event.Group,
		"name":  event.Name,
		"value": event.Value,
	}

	q, _ := c.statements.Get(addEventName)
	_, err := q.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(ctx context.Context, event *model.PageViewHit) error {
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

	q, _ := c.statements.Get(addPageViewName)
	_, err := q.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	paramMap := map[string]interface{}{
		"bid":         event.BID,
		"duration_ms": event.DurationMs,
	}

	q, _ := c.statements.Get(updatePageViewName)
	_, err := q.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
