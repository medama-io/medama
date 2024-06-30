package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
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
			CAST(ifnull(quantile_cont(duration_ms, 0.5), 0) AS INTEGER) AS duration,
			SUM(duration_ms) AS total_duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING duration > 0 AND visitors >= 3`)
	query.WriteString(filter.PaginationString())
	query.WriteString(`)`)
	query.WriteString(`--sql
		SELECT
			pathname,
			duration,
			ifnull(ROUND(total_duration / (SELECT SUM(total_duration) FROM durations), 4), 0) AS duration_percentage
		FROM durations
		ORDER BY duration_percentage DESC, duration DESC, pathname ASC`)

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var time model.StatsTimeSummary
		err := rows.StructScan(&time)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		times = append(times, &time)
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
	query.WriteString(`--sql
		WITH durations AS MATERIALIZED (
		SELECT
			pathname,
			CAST(ifnull(quantile_cont(duration_ms, 0.5), 0) AS INTEGER) AS duration,
			CAST(ifnull(quantile_cont(duration_ms, 0.75), 0) AS INTEGER) AS duration_upper_quartile,
			CAST(ifnull(quantile_cont(duration_ms, 0.25), 0) AS INTEGER) AS duration_lower_quartile,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			SUM(duration_ms) AS total_duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING duration > 0 AND duration_lower_quartile > 0 AND visitors >= 3`)
	query.WriteString(filter.PaginationString())
	query.WriteString(`)`)
	query.WriteString(`--sql
		SELECT
			pathname,
			duration,
			duration_upper_quartile,
			duration_lower_quartile,
			ifnull(ROUND(total_duration / (SELECT SUM(total_duration) FROM durations), 4), 0) AS duration_percentage,
			visitors
		FROM durations
		ORDER BY visitors DESC, duration DESC, pathname ASC`)

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var time model.StatsTime
		err := rows.StructScan(&time)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		times = append(times, &time)
	}

	return times, nil
}
