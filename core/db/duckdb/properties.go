package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

// GetWebsiteReferrersSummary returns a summary of the referrers for the given filters.
func (c *Client) GetWebsiteCustomProperties(ctx context.Context, filter *db.Filters) ([]*model.StatsCustomProperties, error) {
	var properties []*model.StatsCustomProperties

	// Array of custom properties
	//
	// Name is the event key name
	//
	// Value is the event value
	//
	// Events is the number of events for the custom property
	//
	// Visitors is the number of unique visitors for the custom property.
	query := qb.New().
		Select(
			"name",
			"value",
			"COUNT(*) FILTER(value IS NOT NULL) AS events",
			VisitorsStmt,
		).
		From("views").
		LeftJoin("events USING (bid)").
		Where(filter.WhereString()).
		GroupBy("name", "value").
		OrderBy("events DESC", "name ASC", "value ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	defer rows.Close()

	for rows.Next() {
		var property model.StatsCustomProperties
		err := rows.StructScan(&property)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		properties = append(properties, &property)
	}

	return properties, nil
}