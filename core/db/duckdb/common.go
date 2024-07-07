package duckdb

import qb "github.com/medama-io/medama/db/duckdb/query"

// These are common query builder clauses.
const (
	// VisitorsStmt is the number of unique visitors for the query.
	VisitorsStmt = "COUNT(*) FILTER (is_unique_user = true) AS visitors"
	// VisitorsPercentageStmt is the percentage the country contributes to the total unique visits.
	// This expects a CTE named total with a total_visitors column to be present.
	VisitorsPercentageStmt = "ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage"
	// PageviewsStmt is the number of pageviews for the query.
	PageviewsStmt = "COUNT(*) AS pageviews"
	// DurationStmt is the median duration of the visits.
	DurationStmt = "CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration"
	// BouncePercentageStmt is the bounce rate of the visits.
	// This expects a CTE named bounces with a bounce_rate column to be present.
	BounceRateStmt = "(SELECT bounce_rate FROM bounces) AS bounce_rate"
)

// TotalVisitorsCTE declares a materialized CTE to calculate the total number of unique visitors.
func TotalVisitorsCTE(whereClause string) qb.CTE {
	return qb.NewCTE("total", qb.New().
		Select("COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors").
		From("views").
		Where(whereClause))
}

// BounceRateCTE declares a materialized CTE to calculate the bounce rate.
func BounceRateCTE(whereClause string) qb.CTE {
	return qb.NewCTE("bounces", qb.New().
		Select(
			"COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms BETWEEN 1 AND 5000) AS count",
			// Slightly modified to ONLY include visitors that have a duration_ms to
			// represent bounce rates more fairly.
			"COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms IS NOT NULL) AS total_visitors",
			// The bounce rate is calculated as bounces / total_bounce_visitors.
			"ifnull(ROUND(count / total_visitors, 4), 0) AS bounce_rate",
		).
		From("views").
		Where(whereClause))
}
