package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteTimeSummary returns a summary of the time for the given hostname.
func (c *Client) GetWebsiteTimeSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsTimeSummary, error) {
	var times []*model.StatsTimeSummary
	var query strings.Builder

	// Array of time summaries
	//
	// Pathname is the pathname of the page.
	//
	// Duration is the median duration the user spent on the page in milliseconds.
	//
	// DurationPercentage is the percentage the pathname contributes to the total duration.
	query.WriteString(`--sql
		SELECT
			pathname,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration,
			ifnull(ROUND(SUM(duration_ms) * 100.0 / (SELECT SUM(duration_ms) FROM views WHERE hostname = ?), 2), 0) AS duration_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY pathname ORDER BY duration DESC;`)

	err := c.SelectContext(ctx, &times, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return times, nil
}

// GetWebsiteTime returns the time for the given hostname.
func (c *Client) GetWebsiteTime(ctx context.Context, filter *db.Filters) ([]*model.StatsTime, error) {
	var times []*model.StatsTime
	var query strings.Builder

	// Array of time summaries
	//
	// Pathname is the pathname of the page.
	//
	// Duration is the median duration the user spent on the page in milliseconds.
	//
	// DurationUpperQuartile is the upper quartile of the duration the user spent on the page.
	//
	// DurationLowerQuartile is the lower quartile of the duration the user spent on the page.
	//
	// DurationPercentage is the percentage the pathname contributes to the total duration.
	//
	// Title is the title of the page.
	//
	// Bounces is the total number of bounces for the page.
	//
	// Uniques is the total number of uniques for the page.
	query.WriteString(`--sql
		SELECT
			pathname,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration,
			CAST(ifnull(quantile_cont(duration_ms, 0.75), 0) AS INTEGER) AS duration_upper_quartile,
			CAST(ifnull(quantile_cont(duration_ms, 0.25), 0) AS INTEGER) AS duration_lower_quartile,
			ifnull(ROUND(SUM(duration_ms) * 100.0 / (SELECT SUM(duration_ms) FROM views WHERE hostname = ?), 2), 0) AS duration_percentage,
			title,
			COUNT(CASE WHEN duration_ms < 5000 THEN 1 END) AS bounces,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY pathname, title ORDER BY duration DESC;`)

	err := c.SelectContext(ctx, &times, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return times, nil
}
