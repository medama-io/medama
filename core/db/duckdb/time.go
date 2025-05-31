package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

const (
	durationPercentageStmt = "ifnull(ROUND(total_duration / (SELECT SUM(total_duration) FROM durations), 4), 0) AS duration_percentage"
)

// GetWebsiteTimeSummary returns a summary of the time for the given hostname.
func (c *Client) GetWebsiteTimeSummary(
	ctx context.Context,
	filter *db.Filters,
) ([]*model.StatsTimeSummary, error) {
	var times []*model.StatsTimeSummary

	// Array of time summaries
	//
	// Pathname is the pathname of the page.
	//
	// Duration is the average duration the user spent on the page in milliseconds.
	//
	// DurationPercentage is the percentage the pathname contributes to the total duration.
	durationsCTE := qb.New().
		Select(
			"pathname",
			"COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors",
			"CAST(ifnull(quantile_cont(duration_ms, 0.5), 0) AS INTEGER) AS duration",
			"SUM(duration_ms) AS total_duration",
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("pathname").
		Having("duration > 0 AND visitors >= 3")

	if filter.IsCustomEvent {
		durationsCTE = durationsCTE.
			LeftJoin(EventsJoinStmt)
	}

	query := qb.New().WithMaterialized(qb.NewCTE("durations", durationsCTE)).
		Select(
			"pathname",
			"duration",
			durationPercentageStmt,
		).
		From("durations").
		OrderBy("duration_percentage DESC", "duration DESC", "visitors DESC", "pathname ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
func (c *Client) GetWebsiteTime(
	ctx context.Context,
	filter *db.Filters,
) ([]*model.StatsTime, error) {
	var times []*model.StatsTime

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
	durationsCTE := qb.New().
		Select(
			"pathname",
			"CAST(ifnull(quantile_cont(duration_ms, 0.5), 0) AS INTEGER) AS duration",
			"CAST(ifnull(quantile_cont(duration_ms, 0.75), 0) AS INTEGER) AS duration_upper_quartile",
			"CAST(ifnull(quantile_cont(duration_ms, 0.25), 0) AS INTEGER) AS duration_lower_quartile",
			"COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors",
			"SUM(duration_ms) AS total_duration",
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("pathname").
		Having("duration > 0 AND duration_lower_quartile > 0 AND visitors >= 3")

	if filter.IsCustomEvent {
		durationsCTE = durationsCTE.
			LeftJoin(EventsJoinStmt)
	}
	query := qb.New().WithMaterialized(qb.NewCTE("durations", durationsCTE)).
		Select(
			"pathname",
			"duration",
			"duration_upper_quartile",
			"duration_lower_quartile",
			durationPercentageStmt,
			"visitors",
		).
		From("durations").
		OrderBy("visitors DESC", "duration DESC", "pathname ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
