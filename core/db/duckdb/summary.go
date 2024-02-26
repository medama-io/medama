package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteSummary returns the summary stats for the given website.
func (c *Client) GetWebsiteSummary(ctx context.Context, filter *db.Filters) (*model.StatsSummarySingle, error) {
	var summary model.StatsSummarySingle
	var query strings.Builder

	// Visitors are determined by the number of is_unique_user values that are true.
	//
	// Pageviews are determined by the total count of page views that match the hostname.
	//
	// Bounces are determined by any pageview with a duration of less than 5 seconds
	// as well as if they are unique. The percentage is calculated client side as
	// the number of bounces divided by the number of unique pageviews.
	//
	// Duration is the median duration of all pageviews. It needs to be casted to an integer as
	// the median function can return a float for an even number of rows.
	//
	// Active is the number of unique visitors that have visited the website in the last 5 minutes.
	query.WriteString(`--sql
		SELECT
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			COUNT(*) AS pageviews,
			COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	err := c.GetContext(ctx, &summary, query.String(), filter.Args()...)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return &summary, nil
}

// GetWebsiteIntervals returns the stats for the given website by intervals.
func (c *Client) GetWebsiteIntervals(ctx context.Context, filter *db.Filters, interval api.GetWebsiteIDSummaryInterval) ([]*model.StatsIntervals, error) {
	var resp []*model.StatsIntervals
	var query strings.Builder

	// Delete period from filter argument as we are using time_bucket to group
	// by intervals and therefore filter.WhereString() will generate eroneus SQL.
	endPeriod := filter.PeriodEnd
	filter.PeriodEnd = ""

	var intervalQuery string
	switch interval {
	case api.GetWebsiteIDSummaryIntervalMinute:
		intervalQuery = "1 MINUTE"
	case api.GetWebsiteIDSummaryIntervalHour:
		intervalQuery = "1 HOUR"
	case api.GetWebsiteIDSummaryIntervalDay:
		intervalQuery = "1 DAY"
	case api.GetWebsiteIDSummaryIntervalWeek:
		intervalQuery = "7 DAYS"
	case api.GetWebsiteIDSummaryIntervalMonth:
		intervalQuery = "1 MONTH"
	}

	// Intervals are determined by the number of pageviews that match the hostname
	// and are grouped by the interval.
	//
	// We use the WITH clause to generate a series of intervals with empty visitor and pageview counts.
	// We then JOIN the intervals with the actual pageview counts to fill in the gaps.
	// This is done to ensure that we have a row for every interval even if there are no pageviews.
	query.WriteString(`--sql
		WITH intervals AS MATERIALIZED (
			SELECT
				generate_series as interval
			FROM
				generate_series(?::TIMESTAMPTZ, ?::TIMESTAMPTZ, ?::INTERVAL)
		),
		stats AS MATERIALIZED (
			SELECT
				time_bucket(?::INTERVAL, date_created, ?::TIMESTAMPTZ) AS interval,
				COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
				COUNT(*) AS pageviews
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
			GROUP BY interval
		)
		SELECT
			intervals.interval::VARCHAR AS interval,
			COALESCE(stats.visitors, 0) AS visitors,
			COALESCE(stats.pageviews, 0) AS pageviews
		FROM intervals LEFT JOIN stats USING (interval)
		ORDER BY interval ASC`)

	err := c.SelectContext(ctx, &resp, query.String(), filter.Args(filter.PeriodStart, endPeriod, intervalQuery, intervalQuery, filter.PeriodStart)...)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return resp, nil
}
