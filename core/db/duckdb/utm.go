package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

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
		SELECT
			utm_source AS source,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors * 100.0 / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 2), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY utm_source ORDER BY visitors DESC;`)

	err := c.SelectContext(ctx, &utms, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
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
		SELECT
			utm_medium AS medium,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors * 100.0 / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 2), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY utm_medium ORDER BY visitors DESC;`)

	err := c.SelectContext(ctx, &utms, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
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
		SELECT
			utm_campaign AS campaign,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors * 100.0 / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 2), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY utm_campaign ORDER BY visitors DESC;`)

	err := c.SelectContext(ctx, &utms, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return utms, nil
}
