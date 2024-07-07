package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

// GetWebsiteReferrersSummary returns a summary of the referrers for the given filters.
func (c *Client) GetWebsiteReferrersSummary(ctx context.Context, isGroup bool, filter *db.Filters) ([]*model.StatsReferrerSummary, error) {
	var referrers []*model.StatsReferrerSummary

	referrerStmt := "referrer_host AS referrer"
	if isGroup {
		referrerStmt = "IF(referrer_group == '', referrer_host, referrer_group) AS referrer"
	}

	// Array of referrer summaries
	//
	// Referrer is the referrer URL. If isGroup is true, the referrer is the grouped aggregation
	// name. e.g. www.google.com --> Google.
	//
	// Visitors is the number of unique visitors for the referrer.
	//
	// VisitorsPercentage is the percentage the referrer contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			referrerStmt,
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("referrer").
		OrderBy("visitors DESC", "referrer ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
func (c *Client) GetWebsiteReferrers(ctx context.Context, isGroup bool, filter *db.Filters) ([]*model.StatsReferrers, error) {
	var referrers []*model.StatsReferrers

	referrerStmt := "referrer_host AS referrer"
	if isGroup {
		referrerStmt = "IF(referrer_group == '', referrer_host, referrer_group) AS referrer"
	}

	// Array of referrers
	//
	// Referrer is the referrer URL. If isGroup is true, the referrer is the grouped aggregation
	// name. e.g. www.google.com --> Google.
	//
	// Visitors is the number of unique visitors for the referrer.
	//
	// VisitorsPercentage is the percentage the referrer contributes to the total visitors.
	//
	// Bounces is the number of bounces for the referrer.
	//
	// Duration is the median duration the user spent on the page in milliseconds.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			referrerStmt,
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("referrer").
		OrderBy("visitors DESC", "referrer ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
