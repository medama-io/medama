package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

func (c *Client) GetWebsiteUTMSourcesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMSourcesSummary, error) {
	var utms []*model.StatsUTMSourcesSummary

	// Array of utm sources
	//
	// Source is the utm source. Ignore if empty.
	//
	// Visitors is the number of unique visitors for the utm source.
	//
	// VisitorsPercentage is the percentage the utm source contributes to the total unique visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_source AS source",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_source").
		OrderBy("visitors DESC", "source ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMSourcesSummary
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}

// GetWebsiteUTMSources returns the utm sources for the given hostname.
func (c *Client) GetWebsiteUTMSources(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMSources, error) {
	var utms []*model.StatsUTMSources

	// Array of utm sources
	//
	// Source is the utm source. Ignore if empty.
	//
	// Visitors is the number of unique visitors for the utm source.
	//
	// VisitorsPercentage is the percentage the utm source contributes to the total unique visitors.
	//
	// Bounces is the number of unique visitors that bounced for the utm source.
	//
	// Duration is the median duration of the unique visitors for the utm source.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_source AS source",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_source").
		OrderBy("visitors DESC", "source ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMSources
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}

func (c *Client) GetWebsiteUTMMediumsSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMMediumsSummary, error) {
	var utms []*model.StatsUTMMediumsSummary

	// Array of utm mediums
	//
	// Medium is the utm medium.
	//
	// Visitors is the number of unique visitors for the utm medium.
	//
	// VisitorsPercentage is the percentage the utm medium contributes to the total unique visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_medium AS medium",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_medium").
		OrderBy("visitors DESC", "medium ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMMediumsSummary
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}

// GetWebsiteUTMMediums returns the utm mediums for the given hostname.
func (c *Client) GetWebsiteUTMMediums(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMMediums, error) {
	var utms []*model.StatsUTMMediums

	// Array of utm mediums
	//
	// Medium is the utm medium.
	//
	// Visitors is the number of unique visitors for the utm medium.
	//
	// VisitorsPercentage is the percentage the utm medium contributes to the total unique visitors.
	//
	// Bounces is the number of unique visitors that bounced for the utm medium.
	//
	// Duration is the median duration of the unique visitors for the utm medium.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_medium AS medium",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_medium").
		OrderBy("visitors DESC", "medium ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMMediums
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}

func (c *Client) GetWebsiteUTMCampaignsSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMCampaignsSummary, error) {
	var utms []*model.StatsUTMCampaignsSummary

	// Array of utm campaigns
	//
	// Campaign is the utm campaign.
	//
	// Visitors is the number of unique visitors for the utm campaign.
	//
	// VisitorsPercentage is the percentage the utm campaign contributes to the total unique visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_campaign AS campaign",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_campaign").
		OrderBy("visitors DESC", "campaign ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMCampaignsSummary
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}

// GetWebsiteUTMCampaigns returns the utm campaigns for the given hostname.
func (c *Client) GetWebsiteUTMCampaigns(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMCampaigns, error) {
	var utms []*model.StatsUTMCampaigns

	// Array of utm campaigns
	//
	// Campaign is the utm campaign.
	//
	// Visitors is the number of unique visitors for the utm campaign.
	//
	// VisitorsPercentage is the percentage the utm campaign contributes to the total unique visitors.
	//
	// Bounces is the number of unique visitors that bounced for the utm campaign.
	//
	// Duration is the median duration of the unique visitors for the utm campaign.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"utm_campaign AS campaign",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("utm_campaign").
		OrderBy("visitors DESC", "campaign ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var utm model.StatsUTMCampaigns
		err := rows.StructScan(&utm)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		utms = append(utms, &utm)
	}

	return utms, nil
}
