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
	// Duration is the average duration the user spent on the page in milliseconds.
	//
	// DurationPercentage is the percentage the pathname contributes to the total duration.
	query.WriteString(`--sql
		WITH durations AS MATERIALIZED (
		SELECT
			pathname,
			CAST(ifnull(AVG(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING duration > 0 ORDER BY duration DESC, pathname ASC`)
	query.WriteString(filter.PaginationString())
	query.WriteString(`)`)
	query.WriteString(`--sql
		SELECT
			pathname,
			duration,
			ifnull(ROUND(duration * 100.0 / (SELECT SUM(duration) FROM durations), 2), 0) AS duration_percentage
		FROM durations
		ORDER BY duration DESC`)
	err := c.SelectContext(ctx, &times, query.String(), filter.Args()...)
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
	// Duration is the average duration the user spent on the page in milliseconds.
	//
	// DurationUpperQuartile is the upper quartile of the duration the user spent on the page.
	//
	// DurationLowerQuartile is the lower quartile of the duration the user spent on the page.
	//
	// DurationPercentage is the percentage the pathname contributes to the total duration.
	//
	// Visitors is the total number of unique visitors for the page.
	//
	// Bounces is the total number of bounces for the page.
	query.WriteString(`--sql
		WITH durations AS MATERIALIZED (
		SELECT
			pathname,
			CAST(ifnull(AVG(duration_ms), 0) AS INTEGER) AS duration,
			CAST(ifnull(quantile_cont(duration_ms, 0.75), 0) AS INTEGER) AS duration_upper_quartile,
			CAST(ifnull(quantile_cont(duration_ms, 0.25), 0) AS INTEGER) AS duration_lower_quartile,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING duration > 0 ORDER BY duration DESC, pathname ASC`)
	query.WriteString(filter.PaginationString())
	query.WriteString(`)`)
	query.WriteString(`--sql
		SELECT
			pathname,
			duration,
			duration_upper_quartile,
			duration_lower_quartile,
			ifnull(ROUND(duration * 100.0 / (SELECT SUM(duration) FROM durations), 2), 0) AS duration_percentage,
			visitors,
			bounces
		FROM durations
		ORDER BY duration DESC`)
	err := c.SelectContext(ctx, &times, query.String(), filter.Args()...)
	if err != nil {
		return nil, err
	}

	return times, nil
}
