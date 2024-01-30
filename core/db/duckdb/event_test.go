package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddPageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageViewHit{
		BID:         "test_bid",
		Hostname:    "medama-test.io",
		Pathname:    "/",
		Referrer:    "medama.io",
		CountryCode: "GB",
		Language:    "en",
		BrowserName: model.FirefoxBrowser,
		OS:          model.WindowsOS,
		DeviceType:  model.DesktopDevice,
		UTMSource:   "test_source",
		UTMMedium:   "test_medium",
		UTMCampaign: "test_campaign",
	}

	err := client.AddPageView(ctx, event)
	assert.NoError(err)
}

func TestUpdatePageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageViewHit{
		BID:         "test_updated_bid",
		Hostname:    "medama-test.io",
		Pathname:    "/",
		Referrer:    "medama.io",
		CountryCode: "GB",
		Language:    "en",
		BrowserName: model.FirefoxBrowser,
		OS:          model.WindowsOS,
		DeviceType:  model.DesktopDevice,
		UTMSource:   "test_source",
		UTMMedium:   "test_medium",
		UTMCampaign: "test_campaign",
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
