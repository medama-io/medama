package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddPageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageView{
		BID:              "test_bid",
		Hostname:         "medama-test.io",
		Pathname:         "/",
		ReferrerHostname: "medama.io",
		ReferrerPathname: "/test",
		CountryCode:      "GB",
		Language:         "en",
		BrowserName:      "Firefox",
		OS:               "Windows",
		DeviceType:       "Desktop",
		UTMSource:        "test_source",
		UTMMedium:        "test_medium",
		UTMCampaign:      "test_campaign",
	}

	err := client.AddPageView(ctx, event)
	assert.NoError(err)
}

func TestUpdatePageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageView{
		BID:              "test_updated_bid",
		Hostname:         "medama-test.io",
		Pathname:         "/",
		ReferrerHostname: "medama.io",
		ReferrerPathname: "/test",
		CountryCode:      "GB",
		Language:         "en",
		BrowserName:      "Firefox", // "Firefox
		OS:               "Windows",
		DeviceType:       "Desktop",
		UTMSource:        "test_source",
		UTMMedium:        "test_medium",
		UTMCampaign:      "test_campaign",
	}

	err := client.AddPageView(ctx, event)
	assert.NoError(err)

	event2 := &model.PageViewUpdate{
		BID:        "test_update_bid",
		DurationMs: 100,
	}

	err = client.UpdatePageView(ctx, event2)
	assert.NoError(err)
}
