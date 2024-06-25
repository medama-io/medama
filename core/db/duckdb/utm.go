package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

func (c *Client) GetWebsiteUTMSourcesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsUTMSourcesSummary, error) {
	var utms []*model.StatsUTMSourcesSummary
	var query strings.Builder

	// Array of utm sources
	//
	// Source is the utm source. Ignore if empty.
	//
	// Visitors is the number of unique visitors for the utm source.
	//
	// VisitorsPercentage is the percentage the utm source contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_source AS source,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_source ORDER BY visitors DESC, source ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
	var query strings.Builder

	// Array of utm sources
	//
	// Source is the utm source. Ignore if empty.
	//
	// Visitors is the number of unique visitors for the utm source.
	//
	// VisitorsPercentage is the percentage the utm source contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_source AS source,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_source ORDER BY visitors DESC, source ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
	var query strings.Builder

	// Array of utm mediums
	//
	// Medium is the utm medium.
	//
	// Visitors is the number of unique visitors for the utm medium.
	//
	// VisitorsPercentage is the percentage the utm medium contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_medium AS medium,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_medium ORDER BY visitors DESC, medium ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
	var query strings.Builder

	// Array of utm mediums
	//
	// Medium is the utm medium.
	//
	// Visitors is the number of unique visitors for the utm medium.
	//
	// VisitorsPercentage is the percentage the utm medium contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_medium AS medium,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_medium ORDER BY visitors DESC, medium ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
	var query strings.Builder

	// Array of utm campaigns
	//
	// Campaign is the utm campaign.
	//
	// Visitors is the number of unique visitors for the utm campaign.
	//
	// VisitorsPercentage is the percentage the utm campaign contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_campaign AS campaign,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_campaign ORDER BY visitors DESC, campaign ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
	var query strings.Builder

	// Array of utm campaigns
	//
	// Campaign is the utm campaign.
	//
	// Visitors is the number of unique visitors for the utm campaign.
	//
	// VisitorsPercentage is the percentage the utm campaign contributes to the total unique visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			utm_campaign AS campaign,
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY utm_campaign ORDER BY visitors DESC, campaign ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
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
