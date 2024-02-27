package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteBrowser returns the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsersSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsBrowsersSummary, error) {
	var browsers []*model.StatsBrowsersSummary
	var query strings.Builder

	// Array of browsers
	//
	// Browser is the browser name associated with the page.
	//
	// Visitors is the number of unique visitors for the browser.
	//
	// VisitorsPercentage is the percentage the browser contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_browser AS browser,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY browser ORDER BY visitors DESC, browser ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var browser model.StatsBrowsersSummary
		err := rows.StructScan(&browser)
		if err != nil {
			return nil, err
		}
		browsers = append(browsers, &browser)
	}

	return browsers, nil
}

func (c *Client) GetWebsiteBrowsers(ctx context.Context, filter *db.Filters) ([]*model.StatsBrowsers, error) {
	var browsers []*model.StatsBrowsers
	var query strings.Builder

	// Array of browsers
	//
	// Browser is the browser name associated with the page.
	//
	// Visitors is the number of unique visitors for the browser.
	//
	// VisitorsPercentage is the percentage the browser contributes to the total visitors.
	//
	// Bounces is the number of unique visitors that match the pathname and have a duration of less than 5 seconds.
	//
	// Duration is the median duration of the page.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_browser AS browser,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY browser ORDER BY visitors DESC, browser ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var browser model.StatsBrowsers
		err := rows.StructScan(&browser)
		if err != nil {
			return nil, err
		}
		browsers = append(browsers, &browser)
	}

	return browsers, nil
}

func (c *Client) GetWebsiteOSSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsOSSummary, error) {
	var os []*model.StatsOSSummary
	var query strings.Builder

	// Array of operating systems
	//
	// OS is the operating system associated with the page.
	//
	// Visitors is the number of unique visitors for the operating system.
	//
	// VisitorsPercentage is the percentage the operating contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_os AS os,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY os ORDER BY visitors DESC, os ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o model.StatsOSSummary
		err := rows.StructScan(&o)
		if err != nil {
			return nil, err
		}
		os = append(os, &o)
	}

	return os, nil
}

// GetWebsiteOS returns the operating systems for the given hostname.
func (c *Client) GetWebsiteOS(ctx context.Context, filter *db.Filters) ([]*model.StatsOS, error) {
	var os []*model.StatsOS
	var query strings.Builder

	// Array of operating systems
	//
	// OS is the operating system associated with the page.
	//
	// Visitors is the number of unique visitors for the operating system.
	//
	// VisitorsPercentage is the percentage the operating contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_os AS os,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY os ORDER BY visitors DESC, os ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o model.StatsOS
		err := rows.StructScan(&o)
		if err != nil {
			return nil, err
		}
		os = append(os, &o)
	}

	return os, nil
}

func (c *Client) GetWebsiteDevicesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsDevicesSummary, error) {
	var devices []*model.StatsDevicesSummary
	var query strings.Builder

	// Array of devices
	//
	// Device is the device type associated with the page.
	//
	// Visitors is the number of unique visitors for the device.
	//
	// VisitorsPercentage is the percentage the device contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_device_type AS device,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY device ORDER BY visitors DESC, device ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device model.StatsDevicesSummary
		err := rows.StructScan(&device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
}

// GetWebsiteDevices returns the devices for the given hostname.
func (c *Client) GetWebsiteDevices(ctx context.Context, filter *db.Filters) ([]*model.StatsDevices, error) {
	var devices []*model.StatsDevices
	var query strings.Builder

	// Array of devices
	//
	// Device is the device type associated with the page.
	//
	// Visitors is the number of unique visitors for the device.
	//
	// VisitorsPercentage is the percentage the device contributes to the total visitors.
	query.WriteString(`--sql
		WITH total AS MATERIALIZED (
			SELECT COUNT(*) FILTER (WHERE is_unique_page = true) AS total_visitors
			FROM views
			WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(`--sql
		)
		SELECT
			ua_device_type AS device,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT total_visitors FROM total), 4), 0) AS visitors_percentage,
			COUNT(*) FILTER (WHERE is_unique_page = true AND duration_ms < 5000) AS bounces,
			CAST(ifnull(median(duration_ms), 0) AS INTEGER) AS duration
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY device ORDER BY visitors DESC, device ASC`)
	query.WriteString(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.String(), filter.Args(nil))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device model.StatsDevices
		err := rows.StructScan(&device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
}
