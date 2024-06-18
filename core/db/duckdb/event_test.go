package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddPageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)
	defer client.Close()

	event := &model.PageViewHit{
		Hostname:     "medama-test.io",
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

	err := client.AddPageView(ctx, event)
	assert.NoError(err)
}

func TestAddPageDuration(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)
	defer client.Close()

	event := &model.PageViewDuration{
		PageViewHit: model.PageViewHit{
			Hostname:     "medama-test.io",
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

	err := client.AddPageDuration(ctx, event)
	assert.NoError(err)
}
