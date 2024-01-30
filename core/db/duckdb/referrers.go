package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteReferrersSummary returns a summary of the referrers for the given filters.
func (c *Client) GetWebsiteReferrersSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsReferrerSummary, error) {
	var referrers []*model.StatsReferrerSummary
	var query strings.Builder

	// Array of referrer summaries
	//
	// Referrer is the hostname of the referrer.
	//
	// Uniques is the number of uniques for the referrer.
	//
	// UniquePercentage is the percentage the referrer contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			referrer,
			COUNT(CASE WHEN is_unique_page = true THEN 1 END) AS uniques,
			ifnull(ROUND(COUNT(CASE WHEN is_unique_page = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique_page = true THEN 1 END) FROM views WHERE hostname = ?), 2), 0) AS unique_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY referrer ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &referrers, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return referrers, nil
}

// GetWebsiteReferrers returns the referrers for the given hostname.
func (c *Client) GetWebsiteReferrers(ctx context.Context, filter *db.Filters) ([]*model.StatsReferrers, error) {
	var referrers []*model.StatsReferrers
	var query strings.Builder

	// Array of referrers
	//
	// Referrer is the referrer URL.
	//
	// Uniques is the number of uniques for the referrer.
	//
	// UniquePercentage is the percentage the referrer contributes to the total uniques.
	//
	// Bounces is the number of bounces for the referrer.
	//
	// Duration is the median duration the user spent on the page in milliseconds.
	query.WriteString(`--sql
		SELECT
			referrer,
			COUNT(CASE WHEN is_unique_page = true THEN 1 END) AS uniques,
			ifnull(ROUND(COUNT(CASE WHEN is_unique_page = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique_page = true THEN 1 END) FROM views WHERE hostname = ?), 2), 0) AS unique_percentage,
			COUNT(CASE WHEN is_unique_page = true AND duration_ms < 5000 THEN 1 END) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY referrer, ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &referrers, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return referrers, nil
}
