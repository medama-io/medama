package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsitePagesSummary returns a summary of the pages for the given hostname.
func (c *Client) GetWebsitePagesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsPagesSummary, error) {
	var pages []*model.StatsPagesSummary
	var query strings.Builder

	// Array of page paths and their relevant counts
	//
	// Pathname is the path of the page. If it is empty, it is the homepage and defaults to "/".
	//
	// Uniques is the number of unique visitors that match the pathname.
	//
	// UniquePercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// This is ordered by the number of unique visitors in descending order.
	query.WriteString(`--sql
		SELECT
			pathname,
			COUNT(CASE WHEN is_unique_page = true THEN 1 END) AS uniques,
			ifnull(ROUND((uniques * 100.0 / (SELECT COUNT(CASE WHEN is_unique_page = true THEN 1 END) FROM views WHERE hostname = ?)), 2), 0) AS unique_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY pathname ORDER BY uniques DESC`)

	err := c.SelectContext(ctx, &pages, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// GetWebsitePages returns the pages statistics for the given hostname.
func (c *Client) GetWebsitePages(ctx context.Context, filter *db.Filters) ([]*model.StatsPages, error) {
	var pages []*model.StatsPages
	var query strings.Builder

	// Array of page paths and their relevant counts
	//
	// Pathname is the path of the page. If it is empty, it is the homepage and defaults to "/".
	//
	// Uniques is the number of unique visitors that match the pathname.
	//
	// UniquePercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// Pageviews is the number of pageviews that match the pathname.
	//
	// Bounces is the number of bounces that match the pathname. This is calculated if the duration
	// of the pageview is less than 5000 milliseconds.
	//
	// Duration is the median duration of the pageview that match the pathname in milliseconds.
	//
	// This is ordered by the number of unique visitors in descending order.
	query.WriteString(`--sql
		SELECT
			pathname,
			COUNT(CASE WHEN is_unique_page = true THEN 1 END) AS uniques,
			ifnull(ROUND((uniques * 100.0 / (SELECT COUNT(CASE WHEN is_unique_page = true THEN 1 END) FROM views WHERE hostname = ?)), 2), 0) AS unique_percentage,
			COUNT(*) AS pageviews,
			COUNT(CASE WHEN duration_ms < 5000 THEN 1 END) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY pathname ORDER BY uniques DESC`)

	err := c.SelectContext(ctx, &pages, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return pages, nil
}
