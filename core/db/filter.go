package db

import (
	"strings"
)

// Filter is a struct that contains all the possible filters that can be applied
// to a query.
type Filter struct {
	Hostname         string
	Pathname         string
	ReferrerHostname string

	// Time Periods (in RFC3339 format YYYY-MM-DD)
	PeriodStart    string
	PeriodEnd      string
	PeriodInterval string
}

// String builds the WHERE query string.
func (f *Filter) String() string {
	var query strings.Builder

	if f.Hostname != "" {
		query.WriteString("hostname = ?")
	}

	if f.Pathname != "" {
		query.WriteString(" AND pathname = ?")
	}

	if f.ReferrerHostname != "" {
		query.WriteString(" AND referrer_hostname = ?")
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
func (f *Filter) Args(startValues ...string) []interface{} {
	// Initialize the args with the start values
	args := []interface{}{}
	for _, v := range startValues {
		args = append(args, v)
	}

	if f.Hostname != "" {
		args = append(args, f.Hostname)
	}

	if f.Pathname != "" {
		args = append(args, f.Pathname)
	}

	if f.ReferrerHostname != "" {
		args = append(args, f.ReferrerHostname)
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
