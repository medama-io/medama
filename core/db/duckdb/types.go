package duckdb

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

// GetWebsiteBrowser returns the browsers for the given hostname.
func (c *Client) GetWebsiteBrowsersSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsBrowsersSummary, error) {
	var browsers []*model.StatsBrowsersSummary

	// Array of browsers
	//
	// Browser is the browser name associated with the page.
	//
	// Visitors is the number of unique visitors for the browser.
	//
	// VisitorsPercentage is the percentage the browser contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_browser AS browser",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("browser").
		OrderBy("visitors DESC", "browser ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var browser model.StatsBrowsersSummary
		err := rows.StructScan(&browser)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		browsers = append(browsers, &browser)
	}

	return browsers, nil
}

func (c *Client) GetWebsiteBrowsers(ctx context.Context, filter *db.Filters) ([]*model.StatsBrowsers, error) {
	var browsers []*model.StatsBrowsers

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
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_browser AS browser",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("browser").
		OrderBy("visitors DESC", "browser ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var browser model.StatsBrowsers
		err := rows.StructScan(&browser)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		browsers = append(browsers, &browser)
	}

	return browsers, nil
}

func (c *Client) GetWebsiteOSSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsOSSummary, error) {
	var os []*model.StatsOSSummary

	// Array of operating systems
	//
	// OS is the operating system associated with the page.
	//
	// Visitors is the number of unique visitors for the operating system.
	//
	// VisitorsPercentage is the percentage the operating contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_os AS os",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("os").
		OrderBy("visitors DESC", "os ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var o model.StatsOSSummary
		err := rows.StructScan(&o)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		os = append(os, &o)
	}

	return os, nil
}

// GetWebsiteOS returns the operating systems for the given hostname.
func (c *Client) GetWebsiteOS(ctx context.Context, filter *db.Filters) ([]*model.StatsOS, error) {
	var os []*model.StatsOS

	// Array of operating systems
	//
	// OS is the operating system associated with the page.
	//
	// Visitors is the number of unique visitors for the operating system.
	//
	// VisitorsPercentage is the percentage the operating contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_os AS os",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("os").
		OrderBy("visitors DESC", "os ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var o model.StatsOS
		err := rows.StructScan(&o)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		os = append(os, &o)
	}

	return os, nil
}

func (c *Client) GetWebsiteDevicesSummary(ctx context.Context, filter *db.Filters) ([]*model.StatsDevicesSummary, error) {
	var devices []*model.StatsDevicesSummary

	// Array of devices
	//
	// Device is the device type associated with the page.
	//
	// Visitors is the number of unique visitors for the device.
	//
	// VisitorsPercentage is the percentage the device contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_device_type AS device",
			VisitorsStmt,
			VisitorsPercentageStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("device").
		OrderBy("visitors DESC", "device ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var device model.StatsDevicesSummary
		err := rows.StructScan(&device)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		devices = append(devices, &device)
	}

	return devices, nil
}

// GetWebsiteDevices returns the devices for the given hostname.
func (c *Client) GetWebsiteDevices(ctx context.Context, filter *db.Filters) ([]*model.StatsDevices, error) {
	var devices []*model.StatsDevices

	// Array of devices
	//
	// Device is the device type associated with the page.
	//
	// Visitors is the number of unique visitors for the device.
	//
	// VisitorsPercentage is the percentage the device contributes to the total visitors.
	query := qb.New().
		WithMaterialized(TotalVisitorsCTE(filter.WhereString())).
		Select(
			"ua_device_type AS device",
			VisitorsStmt,
			VisitorsPercentageStmt,
			BounceRateStmt,
			DurationStmt,
		).
		From("views").
		Where(filter.WhereString()).
		GroupBy("device").
		OrderBy("visitors DESC", "device ASC").
		Pagination(filter.PaginationString())

	rows, err := c.NamedQueryContext(ctx, query.Build(), filter.Args(nil))
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	defer rows.Close()

	for rows.Next() {
		var device model.StatsDevices
		err := rows.StructScan(&device)
		if err != nil {
			return nil, errors.Wrap(err, "db")
		}
		devices = append(devices, &device)
	}

	return devices, nil
}
