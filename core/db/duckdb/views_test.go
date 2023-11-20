package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddPageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageView{
		BID:            "test_bid",
		Hostname:       "medama-test.io",
		Pathname:       "/",
		Referrer:       "https://medama.io",
		Title:          "Medama",
		Timezone:       "Europe/London",
		Language:       "en",
		BrowserName:    "Chrome",
		BrowserVersion: "91",
		OS:             "MacOS",
		DeviceType:     "Desktop",
		ScreenWidth:    1920,
		ScreenHeight:   1080,
		DateCreated:    1,
	}

	err := client.AddPageView(ctx, event)
	assert.NoError(err)
}

func TestUpdatePageView(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	event := &model.PageView{
		BID:            "test_updated_bid",
		Hostname:       "medama-test.io",
		Pathname:       "/",
		Referrer:       "https://medama.io",
		Title:          "Medama",
		Timezone:       "Europe/London",
		Language:       "en",
		BrowserName:    "Chrome",
		BrowserVersion: "91",
		OS:             "MacOS",
		DeviceType:     "Desktop",
		ScreenWidth:    1920,
		ScreenHeight:   1080,
		DateCreated:    1,
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
