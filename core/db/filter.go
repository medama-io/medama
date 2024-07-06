package db

import (
	"strings"

	"github.com/medama-io/medama/api"
)

// FilterField represents a mapping of the filter field to the database column.
type FilterField string

const (
	FilterHostname      FilterField = "hostname"
	FilterPathname      FilterField = "pathname"
	FilterReferrer      FilterField = "referrer_host"
	FilterReferrerGroup FilterField = "referrer_group"
	FilterUTMSource     FilterField = "utm_source"
	FilterUTMMedium     FilterField = "utm_medium"
	FilterUTMCampaign   FilterField = "utm_campaign"
	FilterBrowser       FilterField = "ua_browser"
	FilterOS            FilterField = "ua_os"
	FilterDevice        FilterField = "ua_device_type"
	FilterCountry       FilterField = "country"
	FilterLanguage      FilterField = "language_base"

	// Custom operations not used in the filtering API
	// but used in named queries.
	FilterPeriodStart FilterField = "start_period"
	FilterPeriodEnd   FilterField = "end_period"
	FilterLimit       FilterField = "limit"
	FilterOffset      FilterField = "offset"
)

// FilterOperation represents the possible filter operations.
type FilterOperation string

const (
	FilterEquals        FilterOperation = "eq"
	FilterNotEquals     FilterOperation = "neq"
	FilterContains      FilterOperation = "contains"
	FilterNotContains   FilterOperation = "not_contains"
	FilterStartsWith    FilterOperation = "starts_with"
	FilterNotStartsWith FilterOperation = "not_starts_with"
	FilterEndsWith      FilterOperation = "ends_with"
	FilterNotEndsWith   FilterOperation = "not_ends_with"
	FilterIn            FilterOperation = "in"
	FilterNotIn         FilterOperation = "not_in"
)

// FilterStringToValues converts an api.FilterString to a value and FilterOperation.
func FilterStringToValues(filterString api.FilterString) (string, FilterOperation) {
	switch {
	case filterString.Eq.IsSet():
		return filterString.Eq.Value, FilterEquals
	case filterString.Neq.IsSet():
		return filterString.Neq.Value, FilterNotEquals
	case filterString.Contains.IsSet():
		return filterString.Contains.Value, FilterContains
	case filterString.NotContains.IsSet():
		return filterString.NotContains.Value, FilterNotContains
	case filterString.StartsWith.IsSet():
		return filterString.StartsWith.Value, FilterStartsWith
	case filterString.NotStartsWith.IsSet():
		return filterString.NotStartsWith.Value, FilterNotStartsWith
	case filterString.EndsWith.IsSet():
		return filterString.EndsWith.Value, FilterEndsWith
	case filterString.NotEndsWith.IsSet():
		return filterString.NotEndsWith.Value, FilterNotEndsWith
	case filterString.In.IsSet():
		return filterString.In.Value, FilterIn
	case filterString.NotIn.IsSet():
		return filterString.NotIn.Value, FilterNotIn
	default:
		return "", ""
	}
}

// Filter represents a single filter with a field, value, and operation.
type Filter struct {
	Field     FilterField
	Value     string
	Operation FilterOperation
}

// NewFilter creates a new filter.
func NewFilter(field FilterField, param api.OptFilterString) *Filter {
	var value string
	var operation FilterOperation

	if param.IsSet() {
		value, operation = FilterStringToValues(param.Value)
	} else {
		return nil
	}

	return &Filter{
		Field: field,
		// Convert the value to lowercase to make the filter case-insensitive
		Value:     strings.ToLower(value),
		Operation: operation,
	}
}

// filterOperationMap maps FilterOperation to their string representations.
//
//nolint:exhaustive,gochecknoglobals // TODO: Implement IN and NOT IN
var filterOperationMap = map[FilterOperation]string{
	FilterEquals:        "ILIKE",
	FilterNotEquals:     "NOT ILIKE",
	FilterContains:      "contains",
	FilterNotContains:   "NOT contains",
	FilterStartsWith:    "starts_with",
	FilterNotStartsWith: "NOT starts_with",
	FilterEndsWith:      "ends_with",
	FilterNotEndsWith:   "NOT ends_with",
}

// String returns the string representation of the filter combined with the operation.
func (f Filter) String() string {
	//nolint:exhaustive // TODO: Implement IN and NOT IN
	switch f.Operation {
	case FilterEquals, FilterNotEquals:
		// e.g. "lower(hostname) = :hostname"
		return string(f.Field) + " " + filterOperationMap[f.Operation] + " :" + string(f.Field)
	case FilterContains, FilterNotContains, FilterStartsWith, FilterNotStartsWith, FilterEndsWith, FilterNotEndsWith:
		// e.g. "contains(hostname, :hostname)"
		return filterOperationMap[f.Operation] + "(LOWER(" + string(f.Field) + "), LOWER(:" + string(f.Field) + "))"
	default:
		return ""
	}
}

// Filters is a struct that contains all the possible filters that can be applied to a query.
type Filters struct {
	Hostname      string
	Pathname      *Filter
	Referrer      *Filter
	ReferrerGroup *Filter
	UTMSource     *Filter
	UTMMedium     *Filter
	UTMCampaign   *Filter
	Browser       *Filter
	OS            *Filter
	Device        *Filter
	Country       *Filter
	Language      *Filter

	// Time Periods (in RFC3339 format 2017-07-21T17:32:28Z)
	PeriodStart string
	PeriodEnd   string

	// Pagination
	Limit  int
	Offset int
}

// addCondition appends a condition to the query if the filter has a non-empty value.
func addCondition(query *strings.Builder, filter *Filter) {
	if filter != nil {
		query.WriteString(" AND " + filter.String())
	}
}

// orCondition appends a condition to the query if both filters have non-empty values.
func orCondition(query *strings.Builder, filter *Filter, filter2 *Filter) {
	switch {
	case filter != nil && filter2 != nil:
		query.WriteString(" OR (" + filter.String() + " OR " + filter2.String() + ")")
	case filter != nil:
		addCondition(query, filter)
	case filter2 != nil:
		addCondition(query, filter2)
	}
}

// String builds the WHERE query string.
func (f Filters) WhereString() string {
	var query strings.Builder

	// Build the query string
	query.WriteString("hostname = :hostname")
	addCondition(&query, f.Pathname)
	addCondition(&query, f.Referrer)
	orCondition(&query, f.Referrer, f.ReferrerGroup)
	addCondition(&query, f.UTMSource)
	addCondition(&query, f.UTMMedium)
	addCondition(&query, f.UTMCampaign)
	addCondition(&query, f.Browser)
	addCondition(&query, f.OS)
	addCondition(&query, f.Device)
	addCondition(&query, f.Country)
	addCondition(&query, f.Language)

	// Time period filters
	if f.PeriodStart != "" {
		query.WriteString(" AND date_created >= CAST(:start_period AS TIMESTAMPTZ)")
	}
	if f.PeriodEnd != "" {
		query.WriteString(" AND date_created <= CAST(:end_period AS TIMESTAMPTZ)")
	}

	return query.String()
}

func (f Filters) PaginationString() string {
	var query strings.Builder
	if f.Limit > 0 {
		query.WriteString(" LIMIT :limit")
	}

	if f.Offset > 0 {
		query.WriteString(" OFFSET :offset")
	}

	return query.String()
}

// Args returns the arguments for the WHERE query string.
// We need this function to pass filter values into the parameters
// of the query to prevent SQL injection.
//
// The startValues are the values that are passed in addition to the filters.
func (f Filters) Args(customMap *map[string]interface{}) map[string]interface{} {
	// Initialize the args with the start map
	if customMap == nil {
		customMap = &map[string]interface{}{}
	}
	args := *customMap

	args[string(FilterHostname)] = f.Hostname
	args[string(FilterPeriodStart)] = f.PeriodStart
	args[string(FilterPeriodEnd)] = f.PeriodEnd
	args[string(FilterLimit)] = f.Limit
	args[string(FilterOffset)] = f.Offset

	//nolint:exhaustive // No other fields use filter structs
	filterValues := map[FilterField]*Filter{
		FilterPathname:      f.Pathname,
		FilterReferrer:      f.Referrer,
		FilterReferrerGroup: f.ReferrerGroup,
		FilterUTMSource:     f.UTMSource,
		FilterUTMMedium:     f.UTMMedium,
		FilterUTMCampaign:   f.UTMCampaign,
		FilterBrowser:       f.Browser,
		FilterOS:            f.OS,
		FilterDevice:        f.Device,
		FilterCountry:       f.Country,
		FilterLanguage:      f.Language,
	}

	// Add non-empty filter values to args
	for field, filter := range filterValues {
		if filter != nil {
			args[string(field)] = filter.Value
		}
	}

	return args
}
