package duckdb

import (
	"context"
	"strings"

	"github.com/medama-io/medama/model"
)

// GetWebsiteBrowserSummary returns a summary of the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsersSummary(ctx context.Context, filter Filter) ([]*model.StatsBrowserSummary, error) {
	var browsers []*model.StatsBrowserSummary
	var query strings.Builder

	// Array of browser summaries
	//
	// Browser is the browser ID number of the page.
	//
	// Uniques is the number of uniques for the browser.
	//
	// UniquePercentage is the percentage the browser contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			ua_browser AS browser,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY browser ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &browsers, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return browsers, nil
}

// GetWebsiteBrowser returns the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsers(ctx context.Context, filter Filter) ([]*model.StatsBrowsers, error) {
	var browsers []*model.StatsBrowsers
	var query strings.Builder

	// Array of browsers
	//
	// Browser is the browser ID number of the page.
	//
	// Uniques is the number of uniques for the browser.
	//
	// UniquePercentage is the percentage the browser contributes to the total uniques.
	//
	// Version is the version of the browser.
	query.WriteString(`--sql
		SELECT
			ua_browser AS browser,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage,
			ua_version AS version
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY browser, version ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &browsers, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return browsers, nil
}

// GetWebsiteOS returns the operating systems for the given hostname.
func (c *Client) GetWebsiteOS(ctx context.Context, filter Filter) ([]*model.StatsOS, error) {
	var os []*model.StatsOS
	var query strings.Builder

	// Array of operating systems
	//
	// OS is the operating system ID number of the page.
	//
	// Uniques is the number of uniques for the operating system.
	//
	// UniquePercentage is the percentage the operating system contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			ua_os AS os,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY os ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &os, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return os, nil
}

// GetWebsiteDevices returns the devices for the given hostname.
func (c *Client) GetWebsiteDevices(ctx context.Context, filter Filter) ([]*model.StatsDevices, error) {
	var devices []*model.StatsDevices
	var query strings.Builder

	// Array of devices
	//
	// Device is the device ID number of the page.
	//
	// Uniques is the number of uniques for the device.
	//
	// UniquePercentage is the percentage the device contributes to the total uniques.
	query.WriteString(`--sql
		SELECT
			ua_device_type AS device,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE `)
	query.WriteString(filter.String())
	query.WriteString(` GROUP BY device ORDER BY uniques DESC;`)

	err := c.SelectContext(ctx, &devices, query.String(), filter.Args(filter.Hostname)...)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
