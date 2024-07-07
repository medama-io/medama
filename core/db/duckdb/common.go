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
	// BouncesStmt is the number of bounces for the query.
	BouncesStmt = "COUNT(*) FILTER (WHERE is_unique_user = true AND duration_ms < 5000) AS bounces"
	// DurationStmt is the median duration of the visits.
	DurationStmt = "CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration"
)

// TotalVisitorsCTE declares a materialized CTE to calculate the total number of unique visitors.
func TotalVisitorsCTE(whereClause string) *qb.QueryBuilder {
	return qb.New().WithMaterialized("total", qb.New().
		Select("COUNT(*) FILTER (WHERE is_unique_user = true) AS total_visitors").
		From("views").
		Where(whereClause),
	)
}
