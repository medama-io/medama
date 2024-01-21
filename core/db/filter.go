package db

import (
	"strconv"
	"strings"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

type FilterField string

const (
	FilterPathname         FilterField = "pathname"
	FilterReferrerHostname FilterField = "referrer_hostname"
	FilterBrowser          FilterField = "ua_browser"
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

			// Convert value to an enum integer if needed (e.g. browser name)
			switch field {
			case FilterBrowser:
				{
					value = strconv.Itoa(int(model.NewBrowserName(value)))
				}
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
func (f *Filter) String() string {
	switch f.Operation {
	case FilterEquals:
		return string(f.Field) + " = ?"
	case FilterNotEquals:
		return string(f.Field) + " != ?"
	case FilterContains:
		return "contains(" + string(f.Field) + ", ?)"
	case FilterNotContains:
		return "NOT contains(" + string(f.Field) + ", ?)"
	case FilterStartsWith:
		return "starts_with(" + string(f.Field) + ", ?)"
	case FilterNotStartsWith:
		return "NOT starts_with(" + string(f.Field) + ", ?)"
	case FilterEndsWith:
		return "ends_with(" + string(f.Field) + ", ?)"
	case FilterNotEndsWith:
		return "NOT ends_with(" + string(f.Field) + ", ?)"
	// TODO: Implement IN and NOT IN
	default:
		return ""
	}
}

// Filters is a struct that contains all the possible filters that can be applied to a query.
type Filters struct {
	Hostname         string
	Pathname         *Filter
	ReferrerHostname *Filter
	Browser          *Filter

	// Time Periods (in RFC3339 format YYYY-MM-DD)
	PeriodStart    string
	PeriodEnd      string
	PeriodInterval string
}

// String builds the WHERE query string.
func (f *Filters) String() string {
	var query strings.Builder

	// Build the query string
	if f.Hostname != "" {
		query.WriteString("hostname = ?")
	}

	if f.Pathname != nil {
		query.WriteString(" AND " + f.Pathname.String())
	}

	if f.ReferrerHostname != nil {
		query.WriteString(" AND " + f.ReferrerHostname.String())
	}

	if f.Browser != nil {
		query.WriteString(" AND " + f.Browser.String())
	}

	// Time period filters
	if f.PeriodStart != "" {
		query.WriteString(" AND date_created >= strptime(?, '%Y-%m-%d')")
	}

	if f.PeriodEnd != "" {
		query.WriteString(" AND date_created <= strptime(?, '%Y-%m-%d')")
	}

	return query.String()
}

// Args returns the arguments for the WHERE query string.
// We need this function to pass filter values into the parameters
// of the query to prevent SQL injection.
func (f *Filters) Args(startValues ...string) []interface{} {
	// Initialize the args with the start values
	args := []interface{}{}
	for _, v := range startValues {
		args = append(args, v)
	}

	if f.Hostname != "" {
		args = append(args, f.Hostname)
	}

	if f.Pathname != nil {
		args = append(args, f.Pathname.Value)
	}

	if f.ReferrerHostname != nil {
		args = append(args, f.ReferrerHostname.Value)
	}

	if f.Browser != nil {
		args = append(args, f.Browser.Value)
	}

	// Time period filters
	if f.PeriodStart != "" {
		args = append(args, f.PeriodStart)
	}

	if f.PeriodEnd != "" {
		args = append(args, f.PeriodEnd)
	}

	return args
}
