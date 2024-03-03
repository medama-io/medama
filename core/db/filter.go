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

// String returns the string representation of the filter combined with the operation.
func (f Filter) String() string {
	switch f.Operation {
	case FilterEquals:
		return string(f.Field) + " = :" + string(f.Field)
	case FilterNotEquals:
		return string(f.Field) + " != :" + string(f.Field)
	case FilterContains:
		return "contains(" + string(f.Field) + ", :" + string(f.Field) + ")"
	case FilterNotContains:
		return "NOT contains(" + string(f.Field) + ", :" + string(f.Field) + ")"
	case FilterStartsWith:
		return "starts_with(" + string(f.Field) + ", :" + string(f.Field) + ")"
	case FilterNotStartsWith:
		return "NOT starts_with(" + string(f.Field) + ", :" + string(f.Field) + ")"
	case FilterEndsWith:
		return "ends_with(" + string(f.Field) + ", :" + string(f.Field) + ")"
	case FilterNotEndsWith:
		return "NOT ends_with(" + string(f.Field) + ", :" + string(f.Field) + ")"
	// TODO: Implement IN and NOT IN
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

// String builds the WHERE query string.
func (f Filters) WhereString() string {
	var query strings.Builder

	// Build the query string
	query.WriteString("hostname = :hostname")

	if f.Pathname != nil {
		query.WriteString(" AND " + f.Pathname.String())
	}

	if f.Referrer != nil {
		query.WriteString(" AND " + f.Referrer.String())
	}

	if f.UTMSource != nil {
		query.WriteString(" AND " + f.UTMSource.String())
	}

	if f.UTMMedium != nil {
		query.WriteString(" AND " + f.UTMMedium.String())
	}

	if f.UTMCampaign != nil {
		query.WriteString(" AND " + f.UTMCampaign.String())
	}

	if f.Browser != nil {
		query.WriteString(" AND " + f.Browser.String())
	}

	if f.OS != nil {
		query.WriteString(" AND " + f.OS.String())
	}

	if f.Device != nil {
		query.WriteString(" AND " + f.Device.String())
	}

	if f.Country != nil {
		query.WriteString(" AND " + f.Country.String())
	}

	if f.Language != nil {
		query.WriteString(" AND " + f.Language.String())
	}

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

	// Add the filters to the args
	if f.Hostname != "" {
		args[string(FilterHostname)] = f.Hostname
	}

	if f.Pathname != nil {
		args[string(FilterPathname)] = f.Pathname.Value
	}

	if f.Referrer != nil {
		args[string(FilterReferrer)] = f.Referrer.Value
	}

	if f.UTMSource != nil {
		args[string(FilterUTMSource)] = f.UTMSource.Value
	}

	if f.UTMMedium != nil {
		args[string(FilterUTMMedium)] = f.UTMMedium.Value
	}

	if f.UTMCampaign != nil {
		args[string(FilterUTMCampaign)] = f.UTMCampaign.Value
	}

	if f.Browser != nil {
		args[string(FilterBrowser)] = f.Browser.Value
	}

	if f.OS != nil {
		args[string(FilterOS)] = f.OS.Value
	}

	if f.Device != nil {
		args[string(FilterDevice)] = f.Device.Value
	}

	if f.Country != nil {
		args[string(FilterCountry)] = f.Country.Value
	}

	if f.Language != nil {
		args[string(FilterLanguage)] = f.Language.Value
	}

	// Time period filters
	if f.PeriodStart != "" {
		args[string(FilterPeriodStart)] = f.PeriodStart
	}

	if f.PeriodEnd != "" {
		args[string(FilterPeriodEnd)] = f.PeriodEnd
	}

	if f.Limit > 0 {
		args[string(FilterLimit)] = f.Limit
	}

	if f.Offset > 0 {
		args[string(FilterOffset)] = f.Offset
	}

	return args
}
