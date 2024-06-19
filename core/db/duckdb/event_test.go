package duckdb_test

import (
	"github.com/medama-io/medama/model"
)

func testAddPageView(db *duckdbTest) {
	rows := db.client.DB.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")
	var count int
	err := rows.Scan(&count)
	db.assert.NoError(err)
	db.assert.Equal(0, count)

	event := &model.PageViewHit{
		Hostname:     "add-page-view-test.io",
		Pathname:     "/",
		IsUniqueUser: true,
		IsUniquePage: true,
		Referrer:     "medama.io",
		CountryCode:  "GB",
		Language:     "en",
		BrowserName:  model.FirefoxBrowser,
		OS:           model.WindowsOS,
		DeviceType:   model.DesktopDevice,
		UTMSource:    "test_source",
		UTMMedium:    "test_medium",
		UTMCampaign:  "test_campaign",
	}

	err = db.client.AddPageView(db.ctx, event)
	db.assert.NoError(err)

	rows = db.client.DB.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")
	err = rows.Scan(&count)
	db.assert.NoError(err)
	db.assert.Equal(1, count)
}

func testAddPageDuration(db *duckdbTest) {
	rows := db.client.DB.QueryRow("SELECT COUNT(*) FROM duration WHERE hostname = 'add-page-duration-test.io'")
	var count int
	err := rows.Scan(&count)
	db.assert.NoError(err)
	db.assert.Equal(0, count)

	event := &model.PageViewDuration{
		PageViewHit: model.PageViewHit{
			Hostname:     "add-page-duration-test.io",
			Pathname:     "/",
			IsUniqueUser: false,
			IsUniquePage: true,
			Referrer:     "medama.io",
			CountryCode:  "GB",
			Language:     "en",
			BrowserName:  model.FirefoxBrowser,
			OS:           model.WindowsOS,
			DeviceType:   model.DesktopDevice,
			UTMSource:    "test_source",
			UTMMedium:    "test_medium",
			UTMCampaign:  "test_campaign",
		},
		DurationMs: 100,
	}

	err = db.client.AddPageDuration(db.ctx, event)
	db.assert.NoError(err)

	rows = db.client.DB.QueryRow("SELECT COUNT(*) FROM duration WHERE hostname = 'add-page-duration-test.io'")
	err = rows.Scan(&count)
	db.assert.NoError(err)
	db.assert.Equal(1, count)
}
