package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

// GetWebsiteBrowser returns the browsers for the given hostname.
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
	query.WriteString(`--sql
		SELECT
			ua_browser AS browser,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY browser ORDER BY visitors DESC, browser ASC`)
	query.WriteString(filter.PaginationString())

	err := c.SelectContext(ctx, &browsers, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return browsers, nil
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
		SELECT
			ua_os AS os,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY os ORDER BY visitors DESC, os ASC`)
	query.WriteString(filter.PaginationString())

	err := c.SelectContext(ctx, &os, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return os, nil
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
		SELECT
			ua_device_type AS device,
			COUNT(*) FILTER (WHERE is_unique_page = true) AS visitors,
			ifnull(ROUND(visitors / (SELECT COUNT(*) FILTER (WHERE is_unique_page = true) FROM views WHERE hostname = ?), 4), 0) AS visitors_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.WhereString())
	query.WriteString(` GROUP BY device ORDER BY visitors DESC, device ASC`)
	query.WriteString(filter.PaginationString())

	err := c.SelectContext(ctx, &devices, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
