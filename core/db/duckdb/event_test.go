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
		BID:          "test_bid",
		Hostname:     "add-page-view-test.io",
		Pathname:     "/",
		IsUniqueUser: true,
		IsUniquePage: true,
		ReferrerHost: "medama.io",
		CountryCode:  "GB",
		LanguageBase: "English",
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

func TestUpdatePageView(t *testing.T) {
	assert, _, ctx, client := SetupDatabase(t)

	event := &model.PageViewHit{
		BID:          "test_updated_bid",
		Hostname:     "medama-test.io",
		Pathname:     "/",
		ReferrerHost: "medama.io",
		CountryCode:  "GB",
		LanguageBase: "English",
		BrowserName:  model.FirefoxBrowser,
		OS:           model.WindowsOS,
		DeviceType:   model.DesktopDevice,
		UTMSource:    "test_source",
		UTMMedium:    "test_medium",
		UTMCampaign:  "test_campaign",
	}

	err := client.AddPageView(ctx, event)
	assert.NoError(err)

	event2 := &model.PageViewDuration{
		BID:        "test_update_bid",
		DurationMs: 100,
	}

	err = client.UpdatePageView(ctx, event2)
	assert.NoError(err)
}
