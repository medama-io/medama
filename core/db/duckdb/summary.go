package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

// GetWebsiteSummary returns the summary stats for the given website.
func (c *Client) GetWebsiteSummary(ctx context.Context, hostname string) (*model.StatsSummary, error) {
	var summary model.StatsSummary

	// Uniques are determined by the number of is_unique values that are true.
	//
	// Pageviews are determined by the number of rows.
	//
	// Bounces are determined by any pageview with a duration of less than 5 seconds
	// as well as if they are unique. The percentage is calculated client side as
	// the number of bounces divided by the number of unique pageviews.
	//
	// Duration is the median duration of all pageviews. It needs to be casted to an integer as
	// the median function can return a float for an even number of rows.
	//
	// Active is the number of unique visitors that have visited the website in the last 5 minutes.
	exec := `--sql
		SELECT
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			COUNT(*) AS pageviews,
			COUNT(CASE WHEN is_unique = true AND duration_ms < 5000 THEN 1 END) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration,
			COUNT(CASE WHEN is_unique = true AND (date_diff('minute', now(), date_updated) < 5) THEN 1 END) AS active
		FROM views
		WHERE hostname = ?`
	err := c.GetContext(ctx, &summary, exec, hostname)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
