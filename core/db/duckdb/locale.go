package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/model"
)

// GetWebsiteCountries returns the countries for the given hostname.
func (c *Client) GetWebsiteCountries(ctx context.Context, filter Filter) ([]*model.StatsCountries, error) {
	var countries []*model.StatsCountries
	var query strings.Builder

	// Array of countries
	//
	// Country is the country of the visitor.
	//
	// Uniques is the number of uniques for the country.
	//
	// UniquePercentage is the percentage the country contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			country_code AS country,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ifnull(ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2), 0) AS unique_percentage,
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY country ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &countries, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

// GetWebsiteLanguages returns the languages for the given hostname.
func (c *Client) GetWebsiteLanguages(ctx context.Context, filter Filter) ([]*model.StatsLanguages, error) {
	var languages []*model.StatsLanguages
	var query strings.Builder

	// Array of languages
	//
	// Language is the language of the visitor.
	//
	// Uniques is the number of uniques for the language.
	//
	// UniquePercentage is the percentage the language contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			language,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ifnull(ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2), 0) AS unique_percentage,
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY language ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &languages, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return languages, nil
}
