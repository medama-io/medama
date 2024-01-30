package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteSummary returns the summary stats for the given website.
func (c *Client) GetWebsiteSummary(ctx context.Context, filter *db.Filters) (*model.StatsSummarySingle, error) {
	var summary model.StatsSummarySingle
	var query strings.Builder

	// Visitors are determined by the number of is_unique_user values that are true.
	//
	// Pageviews are determined by the total count of page views that match the hostname.
	//
	// Bounces are determined by any pageview with a duration of less than 5 seconds
	// as well as if they are unique. The percentage is calculated client side as
	// the number of bounces divided by the number of unique pageviews.
	//
	// Duration is the median duration of all pageviews. It needs to be casted to an integer as
	// the median function can return a float for an even number of rows.
	//
	// Active is the number of unique visitors that have visited the website in the last 5 minutes.
	query.WriteString(`--sql
		SELECT
			COUNT(*) FILTER (WHERE is_unique_user = true) AS visitors,
			COUNT(*) AS pageviews,
			COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration,
			COUNT(*) FILTER (WHERE is_unique_user = true AND (date_diff('minute', now(), date_created) < 5)) AS active
		FROM views
		WHERE `)
	query.WriteString(filter.String())

	err := c.GetContext(ctx, &summary, query.String(), filter.Args()...)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
