package duckdb

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteCountries returns the countries for the given hostname.
func (c *Client) GetWebsiteCountries(ctx context.Context, filter *db.Filters) ([]*model.StatsCountries, error) {
	var countries []*model.StatsCountries
	var query strings.Builder

	// Array of countries
	//
	// Country is the country code of the visitor.
	//
	// Visitors is the number of unique visitors from the country.
	//
	// VisitorsPercentage is the percentage the country contributes to the total unique visits.
	query.WriteString(`--sql
		SELECT
			country_code AS country,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = :hostname), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY country ORDER BY visitors DESC, country ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var country model.StatsCountries
		err := rows.StructScan(&country)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// GetWebsiteLanguages returns the languages for the given hostname.
func (c *Client) GetWebsiteLanguages(ctx context.Context, filter *db.Filters) ([]*model.StatsLanguages, error) {
	var languages []*model.StatsLanguages
	var query strings.Builder

	// Array of languages
	//
	// Language is the language of the visitor.
	//
	// Visitors is the number of unique visitors for the language.
	//
	// VisitorsPercentage is the percentage the language contributes to the total unique visitors.
	query.WriteString(`--sql
		SELECT
			language,
			COUNT(*) FILTER (is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = :hostname), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY language ORDER BY visitors DESC, language ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var language model.StatsLanguages
		err := rows.StructScan(&language)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		languages = append(languages, &language)
	}

	return languages, nil
}
