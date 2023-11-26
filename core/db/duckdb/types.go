package duckdb

import (
	"context"

	"github.com/medama-io/medama/model"
)

// GetWebsiteBrowserSummary returns a summary of the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsersSummary(ctx context.Context, hostname string) ([]*model.StatsBrowserSummary, error) {
	var browsers []*model.StatsBrowserSummary

	// Array of browser summaries
	//
	// Browser is the browser ID number of the page.
	//
	// Uniques is the number of uniques for the browser.
	//
	// UniquePercentage is the percentage the browser contributes to the total uniques.
	query := `
		SELECT
			ua_browser AS browser,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY browser
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &browsers, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return browsers, nil
}

// GetWebsiteBrowser returns the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsers(ctx context.Context, hostname string) ([]*model.StatsBrowsers, error) {
	var browsers []*model.StatsBrowsers

	// Array of browsers
	//
	// Browser is the browser ID number of the page.
	//
	// Uniques is the number of uniques for the browser.
	//
	// UniquePercentage is the percentage the browser contributes to the total uniques.
	//
	// Version is the version of the browser.
	query := `
		SELECT
			ua_browser AS browser,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage,
			ua_version AS version
		FROM views
		WHERE hostname = ?
		GROUP BY browser, version
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &browsers, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return browsers, nil
}

// GetWebsiteOS returns the operating systems for the given hostname.
func (c *Client) GetWebsiteOS(ctx context.Context, hostname string) ([]*model.StatsOS, error) {
	var oss []*model.StatsOS

	// Array of operating systems
	//
	// OS is the operating system ID number of the page.
	//
	// Uniques is the number of uniques for the operating system.
	//
	// UniquePercentage is the percentage the operating system contributes to the total uniques.
	query := `
		SELECT
			ua_os AS os,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY os
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &oss, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return oss, nil
}

// GetWebsiteDevices returns the devices for the given hostname.
func (c *Client) GetWebsiteDevices(ctx context.Context, hostname string) ([]*model.StatsDevices, error) {
	var devices []*model.StatsDevices

	// Array of devices
	//
	// Device is the device ID number of the page.
	//
	// Uniques is the number of uniques for the device.
	//
	// UniquePercentage is the percentage the device contributes to the total uniques.
	query := `
		SELECT
			ua_device_type AS device,
			COUNT(CASE WHEN is_unique = true THEN 1 END) AS uniques,
			ROUND(COUNT(CASE WHEN is_unique = true THEN 1 END) * 100.0 / (SELECT COUNT(CASE WHEN is_unique = true THEN 1 END) FROM views WHERE hostname = ?), 2) AS unique_percentage
		FROM views
		WHERE hostname = ?
		GROUP BY device
		ORDER BY uniques DESC;`

	err := c.SelectContext(ctx, &devices, query, hostname, hostname)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
