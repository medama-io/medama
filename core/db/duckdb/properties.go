package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

const (
	EventsCountStmt      = "COUNT(*) AS events"
	EventsPercentageStmt = "ifnull(ROUND(events / (SELECT total_events FROM total), 4), 0) AS events_percentage"
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
	// Events percentage is the percentage of events for the custom property
	query := qb.New().WithMaterialized(
		qb.NewCTE("total", qb.New().
			Select("COUNT(*) AS total_events").
			From("events").
			Where(filter.WhereString())))

	// If the property name is empty, return only the property names with their
	// aggregated events and visitors. No values.
	if filter.PropertyName == nil || filter.PropertyName.Value == "" {
		query = query.Select(
			"name",
			"'' AS value",
			EventsCountStmt,
			EventsPercentageStmt,
		).
			From("views").
			LeftJoin(EventsJoinStmt).
			Where(filter.WhereString()).
			GroupBy("name").
			OrderBy("events DESC", "name ASC").
			Pagination(filter.PaginationString())
	} else {
		// If the property name is not empty, return the property name with its
		// values, events and visitors.
		query = query.Select(
			"'' AS name",
			"value",
			EventsCountStmt,
			EventsPercentageStmt,
		).
			From("views").
			LeftJoin(EventsJoinStmt).
			Where(filter.WhereString()).
			GroupBy("name", "value").
			OrderBy("events DESC", "name ASC", "value ASC").
			Pagination(filter.PaginationString())
	}

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
