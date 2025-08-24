package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

// GetWebsiteSummary returns the summary stats for the given website.
func (c *Client) GetWebsiteSummary(
	ctx context.Context,
	filter *db.Filters,
) (*model.StatsSummarySingle, error) {
	var summary model.StatsSummarySingle

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
	query := qb.New().
		Select(
			VisitorsStmt,
			PageviewsStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString())

	if filter.IsCustomEvent {
		query = query.
			LeftJoin(EventsJoinStmt)
	}

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&summary)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
	}

	return &summary, nil
}

// GetWebsiteIntervals returns the stats for the given website by intervals.
func (c *Client) GetWebsiteIntervals(
	ctx context.Context,
	filter *db.Filters,
	interval api.GetWebsiteIDSummaryInterval,
) ([]*model.StatsIntervals, error) {
	var resp []*model.StatsIntervals

	isMonthly := false

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
		isMonthly = true
	}

	// Intervals are determined by the number of pageviews that match the hostname
	// and are grouped by the interval.
	//
	// We use the WITH clause to generate a series of intervals with empty visitor and pageview counts.
	// We then JOIN the intervals with the actual pageview counts to fill in the gaps.
	// This is done to ensure that we have a row for every interval even if there are no pageviews.
	//
	// The monthly interval has to be handled differently as we need to use date_trunc to properly
	// handle month boundaries that have varying days.
	var query *qb.QueryBuilder
	if isMonthly {
		statsQuery := qb.New().
			Select(
				"date_trunc('month', views.date_created) AS interval",
				VisitorsStmt,
				PageviewsStmt,
				BounceRateStmt,
				DurationStmt,
			).
			From("views").
			Where(filter.WhereString()).
			GroupBy("interval")

		if filter.IsCustomEvent {
			statsQuery = statsQuery.
				LeftJoin(EventsJoinStmt)
		}

		query = qb.New().
			WithMaterialized(
				qb.NewCTE("intervals", qb.New().
					Select("generate_series as interval").
					From("generate_series(date_trunc('month', CAST(:start_period AS TIMESTAMPTZ)), date_trunc('month', CAST(:end_period AS TIMESTAMPTZ)), INTERVAL '1 MONTH')"),
				),
			).
			WithMaterialized(qb.NewCTE("stats", statsQuery))
	} else {
		statsQuery := qb.New().
			Select(
				"time_bucket(CAST(:interval_query AS INTERVAL), views.date_created, CAST(:start_period AS TIMESTAMPTZ)) AS interval",
				VisitorsStmt,
				PageviewsStmt,
				BounceRateStmt,
				DurationStmt,
			).
			From("views").
			Where(filter.WhereString()).
			GroupBy("interval")

		if filter.IsCustomEvent {
			statsQuery = statsQuery.
				LeftJoin(EventsJoinStmt)
		}

		query = qb.New().
			WithMaterialized(
				qb.NewCTE("intervals", qb.New().
					Select("generate_series as interval").
					From("generate_series(CAST(:start_period AS TIMESTAMPTZ), CAST(:end_period AS TIMESTAMPTZ), CAST(:interval_query AS INTERVAL))"),
				),
			).
			WithMaterialized(qb.NewCTE("stats", statsQuery))
	}

	query = query.
		Select(
			"CAST(intervals.interval AS VARCHAR) AS interval",
			"COALESCE(stats.visitors, 0) AS visitors",
			"COALESCE(stats.pageviews, 0) AS pageviews",
			"COALESCE(stats.bounce_rate, 0.0) AS bounce_rate",
			"COALESCE(stats.duration, 0) AS duration",
		).
		From("intervals").
		LeftJoin("stats USING (interval)").
		OrderBy("interval ASC")

	filterMap := map[string]any{
		"interval_query": intervalQuery,
	}

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(&filterMap))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var interval model.StatsIntervals

		err := rows.StructScan(&interval)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}

		resp = append(resp, &interval)
	}

	return resp, nil
}

// GetWebsiteSummaryLast24Hours returns the summary stats for the given website in the last 24 hours.
func (c *Client) GetWebsiteSummaryLast24Hours(
	ctx context.Context,
	hostname string,
) (*model.StatsSummaryLast24Hours, error) {
	var summary model.StatsSummaryLast24Hours
	// Visitors are determined by the number of is_unique_user values that are true.
	query := qb.New().
		Select(
			VisitorsStmt,
		).
		From("views").
		Where("hostname = :hostname AND date_created > now() - INTERVAL '1 DAY'")

	filterMap := map[string]any{
		"hostname": hostname,
	}

	rows, err := c.NamedQueryContext(ctx, query.Build(), filterMap)
	if err != nil {
		return nil, errors.Wrap(err, "db: GetWebsiteSummaryLast24Hours")
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&summary)
		if err != nil {
			return nil, errors.Wrap(err, "db: GetWebsiteSummaryLast24Hours")
		}
	}

	return &summary, nil
}
