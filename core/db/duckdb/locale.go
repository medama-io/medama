package duckdb

import (
	"context"
	"strings"

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
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY country ORDER BY visitors DESC, country ASC`)
	query.WriteString(filter.PaginationString())

	err := c.SelectContext(ctx, &countries, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
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
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY language ORDER BY visitors DESC, language ASC`)
	query.WriteString(filter.PaginationString())

	err := c.SelectContext(ctx, &languages, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return languages, nil
}
