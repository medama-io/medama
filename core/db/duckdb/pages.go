package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

const (
	// PageviewsVisitorStmt is the number of unique visitors for the query using is_unique_page.
	PageviewsVisitorStmt = "COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors"
)

// TotalPageviewsCTE declares a materialized CTE to calculate the total number of unique visitors
// per page and pageviews.
func totalPageviewsCTE(whereClause string, isCustomEvent bool) qb.CTE {
	query := qb.New().
		Select(
			"COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors",
			"COUNT(*) AS total_pageviews",
		).
		From("views").
		Where(whereClause)

	if isCustomEvent {
		query = query.
			LeftJoin(EventsJoinStmt)
	}

	return qb.NewCTE("total", query)
}

// GetWebsitePagesSummary returns a summary of the pages for the given hostname.
func (c *Client) GetWebsitePagesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsPagesSummary, error) {
	var pages []*model.StatsPagesSummary

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
	query := qb.New().
		WithMaterialized(totalPageviewsCTE(filter.WhereString(), filter.IsCustomEvent)).
		Select(
			"pathname",
			// Different from VisitorsStmt due to is_unique_page
			PageviewsVisitorStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("pathname").
		Having("visitors > 0").
		OrderBy("visitors DESC", "pathname ASC").
		Pagination(filter.PaginationString())

	if filter.IsCustomEvent {
		query = query.
			LeftJoin(EventsJoinStmt)
	}

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
	query := qb.New().
		WithMaterialized(totalPageviewsCTE(filter.WhereString(), filter.IsCustomEvent)).
		Select(
			"pathname",
			// Different from VisitorsStmt due to is_unique_page
			PageviewsVisitorStmt,
			VisitorsPercentageStmt,
			PageviewsStmt,
			"ifnull(ROUND(pageviews / (SELECT total_pageviews FROM total), 4), 0) AS pageviews_percentage",
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("pathname").
		Having("visitors > 0").
		OrderBy("visitors DESC", "pageviews DESC", "pathname ASC").
		Pagination(filter.PaginationString())

	if filter.IsCustomEvent {
		query = query.
			LeftJoin(EventsJoinStmt)
	}

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
