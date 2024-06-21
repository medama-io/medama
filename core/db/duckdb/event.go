package duckdb

import (
	"context"
	"sync"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

var (
	//nolint:gochecknoglobals // Reason: Singleton patterns are typically written like this.
	addOnce sync.Once
	//nolint:gochecknoglobals // Reason: Prepared statements are meant to be global.
	addStmt *sqlx.NamedStmt
	//nolint:gochecknoglobals // Reason: Singleton patterns are typically written like this.
	updateOnce sync.Once
	//nolint:gochecknoglobals // Reason: Prepared statements are meant to be global.
	updateStmt *sqlx.NamedStmt
)

// AddPageView adds a page view to the database.
func (c *Client) AddPageView(ctx context.Context, event *model.PageViewHit) error {
	var err error
	// Prepare named exec once.
	addOnce.Do(func() {
		exec := `--sql
		INSERT INTO views (
			bid,
			hostname,
			pathname,
			is_unique_user,
			is_unique_page,
			referrer_host,
			referrer_group,
			country_code,
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
			:country_code,
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

		addStmt, err = c.DB.PrepareNamedContext(ctx, exec)
		if err != nil {
			log := logger.Get()
			log.Error().Err(err).Msg("failed to create prepared statement for add page view")
			panic("failed to create prepared statement for add page view")
		}
	})

	paramMap := map[string]interface{}{
		"bid":              event.BID,
		"hostname":         event.Hostname,
		"pathname":         event.Pathname,
		"is_unique_user":   event.IsUniqueUser,
		"is_unique_page":   event.IsUniquePage,
		"referrer_host":    event.ReferrerHost,
		"referrer_group":   event.ReferrerGroup,
		"country_code":     event.CountryCode,
		"language_base":    event.LanguageBase,
		"language_dialect": event.LanguageDialect,
		"ua_browser":       event.BrowserName,
		"ua_os":            event.OS,
		"ua_device_type":   event.DeviceType,
		"utm_source":       event.UTMSource,
		"utm_medium":       event.UTMMedium,
		"utm_campaign":     event.UTMCampaign,
	}

	_, err = addStmt.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}

// UpdatePageView updates a page view in the database.
func (c *Client) UpdatePageView(ctx context.Context, event *model.PageViewDuration) error {
	var err error
	// Prepare named exec once.
	updateOnce.Do(func() {
		exec := `--sql
		UPDATE views SET duration_ms = :duration_ms WHERE bid = :bid`

		updateStmt, err = c.DB.PrepareNamedContext(ctx, exec)
		if err != nil {
			log := logger.Get()
			log.Error().Err(err).Msg("failed to create prepared statement for update page view")
			panic(err)
		}
	})

	paramMap := map[string]interface{}{
		"bid":         event.BID,
		"duration_ms": event.DurationMs,
	}

	_, err = updateStmt.ExecContext(ctx, paramMap)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	return nil
}
