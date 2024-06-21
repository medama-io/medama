package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
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
	// Visitors is the number of unique visitors for the referrer.
	//
	// VisitorsPercentage is the percentage the referrer contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			referrer_host,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY referrer_host ORDER BY visitors DESC, referrer_host ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var referrer model.StatsReferrerSummary
		err := rows.StructScan(&referrer)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		referrers = append(referrers, &referrer)
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
	// Visitors is the number of unique visitors for the referrer.
	//
	// VisitorsPercentage is the percentage the referrer contributes to the total visitors.
	//
	// Bounces is the number of bounces for the referrer.
	//
	// Duration is the median duration the user spent on the page in milliseconds.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			referrer,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY referrer_host ORDER BY visitors DESC, referrer_host ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var referrer model.StatsReferrers
		err := rows.StructScan(&referrer)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		referrers = append(referrers, &referrer)
	}

	return referrers, nil
}
