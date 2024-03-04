package db

import (
	"strconv"
	"strings"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

// FilterField represents a mapping of the filter field to the database column.
type FilterField string

const (
	FilterHostname    FilterField = "hostname"
	FilterPathname    FilterField = "pathname"
	FilterReferrer    FilterField = "referrer"
	FilterUTMSource   FilterField = "utm_source"
	FilterUTMMedium   FilterField = "utm_medium"
	FilterUTMCampaign FilterField = "utm_campaign"
	FilterBrowser     FilterField = "ua_browser"
	FilterOS          FilterField = "ua_os"
	FilterDevice      FilterField = "ua_device_type"
	FilterCountry     FilterField = "country_code"
	FilterLanguage    FilterField = "language"

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

// FilterFixedToValues converts an api.FilterFixed to a string and FilterOperation.
func FilterFixedToValues(filterFixed api.FilterFixed) (string, FilterOperation) {
	switch {
	case filterFixed.Eq.IsSet():
		return filterFixed.Eq.Value, FilterEquals
	case filterFixed.Neq.IsSet():
		return filterFixed.Neq.Value, FilterNotEquals
	case filterFixed.In.IsSet():
		return filterFixed.In.Value, FilterIn
	case filterFixed.NotIn.IsSet():
		return filterFixed.NotIn.Value, FilterNotIn
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
func NewFilter(field FilterField, param interface{}) *Filter {
	var value string
	var operation FilterOperation

	switch v := param.(type) {
	case api.OptFilterFixed:
		if v.IsSet() {
			value, operation = FilterFixedToValues(v.Value)

			// Convert the value to the enum integer for the database
			//nolint:exhaustive // No other fields use uint8 enums
			switch field {
			case FilterBrowser:
				value = strconv.Itoa(int(model.NewBrowserName(value)))
			case FilterOS:
				value = strconv.Itoa(int(model.NewOSName(value)))
			case FilterDevice:
				value = strconv.Itoa(int(model.NewDeviceTypeString(value)))
			default:
				// Do nothing
			}
		} else {
			return nil
		}
	case api.OptFilterString:
		if v.IsSet() {
			value, operation = FilterStringToValues(v.Value)
		} else {
			return nil
		}
	default:
		return nil
	}

	return &Filter{
		Field:     field,
		Value:     value,
		Operation: operation,
	}
}

// filterOperationMap maps FilterOperation to their string representations.
//
//nolint:exhaustive,gochecknoglobals // TODO: Implement IN and NOT IN
var filterOperationMap = map[FilterOperation]string{
	FilterEquals:        "=",
	FilterNotEquals:     "!=",
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
		return string(f.Field) + " " + filterOperationMap[f.Operation] + " :" + string(f.Field)
	case FilterContains, FilterNotContains, FilterStartsWith, FilterNotStartsWith, FilterEndsWith, FilterNotEndsWith:
		return filterOperationMap[f.Operation] + "(" + string(f.Field) + ", :" + string(f.Field) + ")"
	default:
		return ""
	}
}

// Filters is a struct that contains all the possible filters that can be applied to a query.
type Filters struct {
	Hostname    string
	Pathname    *Filter
	Referrer    *Filter
	UTMSource   *Filter
	UTMMedium   *Filter
	UTMCampaign *Filter
	Browser     *Filter
	OS          *Filter
	Device      *Filter
	Country     *Filter
	Language    *Filter

	// Time Periods (in RFC3339 format 2017-07-21T17:32:28Z)
	PeriodStart string
	PeriodEnd   string

	// Pagination
	Limit  int
	Offset int
}

// addCondition appends a condition to the query if the filter has a non-empty value.
func addCondition(query *strings.Builder, condition string, value string) {
	if value != "" {
		query.WriteString(" AND " + condition)
	}
}

// String builds the WHERE query string.
func (f Filters) WhereString() string {
	var query strings.Builder

	// Build the query string
	query.WriteString("hostname = :hostname")
	addCondition(&query, f.Pathname.String(), f.Pathname.Value)
	addCondition(&query, f.Referrer.String(), f.Referrer.Value)
	addCondition(&query, f.UTMSource.String(), f.UTMSource.Value)
	addCondition(&query, f.UTMMedium.String(), f.UTMMedium.Value)
	addCondition(&query, f.UTMCampaign.String(), f.UTMCampaign.Value)
	addCondition(&query, f.Browser.String(), f.Browser.Value)
	addCondition(&query, f.OS.String(), f.OS.Value)
	addCondition(&query, f.Device.String(), f.Device.Value)
	addCondition(&query, f.Country.String(), f.Country.Value)
	addCondition(&query, f.Language.String(), f.Language.Value)

	// Time period filters
	addCondition(&query, "date_created >= CAST(:start_period AS TIMESTAMPTZ)", f.PeriodStart)
	addCondition(&query, "date_created <= CAST(:end_period AS TIMESTAMPTZ)", f.PeriodEnd)

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

	filterValues := map[FilterField]interface{}{
		FilterHostname:    f.Hostname,
		FilterPathname:    f.Pathname.Value,
		FilterReferrer:    f.Referrer.Value,
		FilterUTMSource:   f.UTMSource.Value,
		FilterUTMMedium:   f.UTMMedium.Value,
		FilterUTMCampaign: f.UTMCampaign.Value,
		FilterBrowser:     f.Browser.Value,
		FilterOS:          f.OS.Value,
		FilterDevice:      f.Device.Value,
		FilterCountry:     f.Country.Value,
		FilterLanguage:    f.Language.Value,
		FilterPeriodStart: f.PeriodStart,
		FilterPeriodEnd:   f.PeriodEnd,
		FilterLimit:       f.Limit,
		FilterOffset:      f.Offset,
	}

	// Add non-empty filter values to args
	for field, value := range filterValues {
		if value != "" {
			args[string(field)] = value
		}
	}

	return args
}
