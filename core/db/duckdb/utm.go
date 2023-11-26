package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

// GetWebsiteUTMSources returns the utm sources for the given hostname.
func (c *Client) GetWebsiteUTMSources(ctx context.Context, hostname string) ([]*model.StatsUTMSources, error) {
	var utms []*model.StatsUTMSources

	// Array of utm sources
	//
	// Source is the utm source. Ignore if empty.
	//
	// Uniques is the number of uniques for the utm source.
	//
	// UniquePercentage is the percentage the utm source contributes to the total uniques.
	query := `
		SELECT
			utm_source AS source,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY utm_source
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &utms, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return utms, nil
}

// GetWebsiteUTMMediums returns the utm mediums for the given hostname.
func (c *Client) GetWebsiteUTMMediums(ctx context.Context, hostname string) ([]*model.StatsUTMMediums, error) {
	var utms []*model.StatsUTMMediums

	// Array of utm mediums
	//
	// Medium is the utm medium.
	//
	// Uniques is the number of uniques for the utm medium.
	//
	// UniquePercentage is the percentage the utm medium contributes to the total uniques.
	query := `
		SELECT
			utm_medium AS medium,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY utm_medium
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &utms, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return utms, nil
}

// GetWebsiteUTMCampaigns returns the utm campaigns for the given hostname.
func (c *Client) GetWebsiteUTMCampaigns(ctx context.Context, hostname string) ([]*model.StatsUTMCampaigns, error) {
	var utms []*model.StatsUTMCampaigns

	// Array of utm campaigns
	//
	// Campaign is the utm campaign.
	//
	// Uniques is the number of uniques for the utm campaign.
	//
	// UniquePercentage is the percentage the utm campaign contributes to the total uniques.
	query := `
		SELECT
			utm_campaign AS campaign,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY utm_campaign
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &utms, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return utms, nil
}
