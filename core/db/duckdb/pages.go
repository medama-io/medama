package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

// GetWebsitePagesSummary returns a summary of the pages for the given hostname.
func (c *Client) GetWebsitePagesSummary(ctx context.Context, hostname string) ([]*model.StatsPagesSummary, error) {
	var pages []*model.StatsPagesSummary

	// Array of page paths and their relevant counts
	//
	// Pathname is the path of the page. If it is empty, it is the homepage and defaults to "/".
	//
	// Uniques is the number of unique visitors that match the pathname.
	//
	// UniquesPercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// This is ordered by the number of unique visitors in descending order.
	query := `--sql
		SELECT
			pathname,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			(uniques / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?)) * 100 AS uniques_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY pathname
		ORDER BY uniques, pathname DESC`

	err := c.SelectContext(ctx, &pages, query, hostname)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// GetWebsitePages returns the pages statistics for the given hostname.
func (c *Client) GetWebsitePages(ctx context.Context, hostname string) ([]*model.StatsPages, error) {
	var pages []*model.StatsPages

	// Array of page paths and their relevant counts
	//
	// Pathname is the path of the page. If it is empty, it is the homepage and defaults to "/".
	//
	// Uniques is the number of unique visitors that match the pathname.
	//
	// UniquesPercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// Pageviews is the number of pageviews that match the pathname.
	//
	// Bounces is the number of bounces that match the pathname. This is calculated if the duration
	// of the pageview is less than 5000 milliseconds.
	//
	// Duration is the average duration of the pageview that match the pathname in milliseconds.
	//
	// This is ordered by the number of unique visitors in descending order.
	query := `--sql
		SELECT
			pathname,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			(uniques / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?)) * 100 AS uniques_percentage,
			COUNT(*) AS pageviews,
			COUNT(CASE WHEN duration_ms < 5000 THEN 1 END) AS bounces,
			CAST(AVG(duration_ms) AS INTEGER) AS duration
		FROM views
		WHERE hostname = ?
		GROUP BY pathname
		ORDER BY uniques, pageviews, pathname DESC`
	// TODO: Switch to DECIMAL for uniques_percentage

	err := c.SelectContext(ctx, &pages, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return pages, nil
}
