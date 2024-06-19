package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddPageView(t *testing.T) {
	assert, _, ctx, client := SetupDatabase(t)
	rows := client.DB.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")
	var count int
	err := rows.Scan(&count)
	assert.NoError(err)
	assert.Equal(0, count)

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

	err = client.AddPageView(ctx, event)
	assert.NoError(err)

	rows = client.DB.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")
	err = rows.Scan(&count)
	assert.NoError(err)
	assert.Equal(1, count)
}

func testAddPageDuration(t *testing.T) {
	assert, _, ctx, client := SetupDatabase(t)
	rows := client.DB.QueryRow("SELECT COUNT(*) FROM duration WHERE hostname = 'add-page-duration-test.io'")
	var count int
	err := rows.Scan(&count)
	assert.NoError(err)
	assert.Equal(0, count)

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

	err = client.AddPageDuration(ctx, event)
	assert.NoError(err)

	rows = client.DB.QueryRow("SELECT COUNT(*) FROM duration WHERE hostname = 'add-page-duration-test.io'")
	err = rows.Scan(&count)
	assert.NoError(err)
	assert.Equal(1, count)
}
