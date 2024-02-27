package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
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
	// Visitors is the number of unique visitors that match the pathname.
	//
	// VisitorsPercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// This is ordered by the number of unique visitors in descending order.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			pathname,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING visitors > 0 ORDER BY visitors DESC, pathname ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var page model.StatsPagesSummary
		err := rows.StructScan(&page)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		pages = append(pages, &page)
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
	// Visitors is the number of unique visitors that match the pathname.
	//
	// VisitorsPercentage is the percentage of unique visitors that match the pathname
	// out of all unique visitors for the website.
	//
	// Pageviews is the number of pageviews that match the pathname.
	//
	// PageviewsPercentage is the percentage of pageviews that match the pathname out of all
	// pageviews for the website.
	//
	// Bounces is the number of bounces that match the pathname. This is calculated if the duration
	// of the pageview is less than 5000 milliseconds.
	//
	// Duration is the median duration of the pageview that match the pathname in milliseconds.
	//
	// This is ordered by the number of unique visitors in descending order.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT
				COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors,
				COUNT(*) AS total_pageviews
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			pathname,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) AS pageviews,
			ifnull(ROUND(pageviews / (SELECT total_pageviews FROM total), 4), 0) AS pageviews_percentage,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY pathname HAVING visitors > 0 ORDER BY visitors DESC, pageviews DESC, pathname ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var page model.StatsPages
		err := rows.StructScan(&page)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		pages = append(pages, &page)
	}

	return pages, nil
}
