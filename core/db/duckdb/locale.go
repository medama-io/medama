package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

// GetWebsiteCountries returns the countries for the given hostname.
func (c *Client) GetWebsiteCountries(ctx context.Context, hostname string) ([]*model.StatsCountries, error) {
	var countries []*model.StatsCountries

	// Array of countries
	//
	// Country is the country of the visitor.
	//
	// Uniques is the number of uniques for the country.
	//
	// UniquePercentage is the percentage the country contributes to the total uniques.
	query := `
		SELECT
			country_code AS country,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage,
		FROM views
		WHERE hostname = ?
		GROUP BY country
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &countries, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

// GetWebsiteLanguages returns the languages for the given hostname.
func (c *Client) GetWebsiteLanguages(ctx context.Context, hostname string) ([]*model.StatsLanguages, error) {
	var languages []*model.StatsLanguages

	// Array of languages
	//
	// Language is the language of the visitor.
	//
	// Uniques is the number of uniques for the language.
	//
	// UniquePercentage is the percentage the language contributes to the total uniques.
	query := `
		SELECT
			language,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage,
		FROM views
		WHERE hostname = ?
		GROUP BY language
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &languages, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return languages, nil
}
