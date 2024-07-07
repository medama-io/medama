package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

// GetWebsiteCountries returns the countries for the given hostname.
func (c *Client) GetWebsiteCountriesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsCountriesSummary, error) {
	var countries []*model.StatsCountriesSummary

	// Array of countries
	//
	// Country is the country code of the visitor.
	//
	// Visitors is the number of unique visitors from the country.
	//
	// VisitorsPercentage is the percentage the country contributes to the total unique visits.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"country",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("country").
		OrderBy("visitors DESC", "country ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var country model.StatsCountriesSummary
		err := rows.StructScan(&country)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		countries = append(countries, &country)
	}

	return countries, nil
}

// GetWebsiteCountries returns the countries for the given hostname.
func (c *Client) GetWebsiteCountries(ctx context.Context, filter *db.Filters) ([]*model.StatsCountries, error) {
	var countries []*model.StatsCountries

	// Array of countries
	//
	// Country is the country code of the visitor.
	//
	// Visitors is the number of unique visitors from the country.
	//
	// VisitorsPercentage is the percentage the country contributes to the total unique visits.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"country",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("country").
		OrderBy("visitors DESC", "country ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
func (c *Client) GetWebsiteLanguagesSummary(ctx context.Context, isLocale bool, filter *db.Filters) ([]*model.StatsLanguagesSummary, error) {
	var languages []*model.StatsLanguagesSummary

	languageSelect := "language_base AS language"
	if isLocale {
		languageSelect = "language_dialect AS language"
	}

	// Array of languages
	//
	// Language is the language of the visitor. If isLocale is true, the language is the locale of the visitor.
	//
	// Visitors is the number of unique visitors for the language.
	//
	// VisitorsPercentage is the percentage the language contributes to the total unique visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			languageSelect,
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("language").
		OrderBy("visitors DESC", "language ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var language model.StatsLanguagesSummary
		err := rows.StructScan(&language)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		languages = append(languages, &language)
	}

	return languages, nil
}

// GetWebsiteLanguages returns the languages for the given hostname.
func (c *Client) GetWebsiteLanguages(ctx context.Context, isLocale bool, filter *db.Filters) ([]*model.StatsLanguages, error) {
	var languages []*model.StatsLanguages

	languageSelect := "language_base AS language"
	if isLocale {
		languageSelect = "language_dialect AS language"
	}

	// Array of languages
	//
	// Language is the language of the visitor.
	//
	// Visitors is the number of unique visitors for the language.
	//
	// VisitorsPercentage is the percentage the language contributes to the total unique visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			languageSelect,
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("language").
		OrderBy("visitors DESC", "language ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
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
