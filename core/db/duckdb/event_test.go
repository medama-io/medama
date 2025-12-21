package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestAddEvents(t *testing.T) {
	assert, require, ctx, client := SetupDatabase(t)
	rows := client.QueryRow("SELECT COUNT(*) FROM events WHERE group_name = 'add-event-test.io'")

	var count int

	err := rows.Scan(&count)
	require.NoError(err)
	assert.Equal(0, count)

	event1 := model.EventHit{
		Group: "add-event-test.io",
		Name:  "test_event",
		Value: "test_value",
	}

	event2 := model.EventHit{
		Group: "add-event-test.io",
		Name:  "test_event2",
		Value: "test_value2",
	}

	err = client.AddEvents(ctx, &[]model.EventHit{event1, event2})
	require.NoError(err)

	rows = client.QueryRow("SELECT COUNT(*) FROM events WHERE group_name = 'add-event-test.io'")
	err = rows.Scan(&count)
	require.NoError(err)
	assert.Equal(2, count)
}

func TestAddPageView(t *testing.T) {
	assert, require, ctx, client := SetupDatabase(t)
	rows := client.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")

	var count int

	err := rows.Scan(&count)
	require.NoError(err)
	assert.Equal(0, count)

	event := &model.PageViewHit{
		BID:          "test_bid",
		Hostname:     "add-page-view-test.io",
		Pathname:     "/",
		IsUniqueUser: true,
		IsUniquePage: true,
		ReferrerHost: "medama.io",
		Country:      "United Kingdom",
		LanguageBase: "English",
		BrowserName:  "Firefox",
		OS:           "Windows",
		DeviceType:   "Desktop",
		UTMSource:    "test_source",
		UTMMedium:    "test_medium",
		UTMCampaign:  "test_campaign",
	}

	err = client.AddPageView(ctx, event, nil)
	require.NoError(err)

	rows = client.QueryRow("SELECT COUNT(*) FROM views WHERE hostname = 'add-page-view-test.io'")
	err = rows.Scan(&count)
	require.NoError(err)
	assert.Equal(1, count)
}

func TestUpdatePageView(t *testing.T) {
	assert, require, ctx, client := SetupDatabase(t)

	event := &model.PageViewHit{
		BID:          "test_updated_bid",
		Hostname:     "medama-test.io",
		Pathname:     "/",
		ReferrerHost: "medama.io",
		Country:      "United Kingdom",
		LanguageBase: "English",
		BrowserName:  "Firefox",
		OS:           "Windows",
		DeviceType:   "Desktop",
		UTMSource:    "test_source",
		UTMMedium:    "test_medium",
		UTMCampaign:  "test_campaign",
	}

	err := client.AddPageView(ctx, event, nil)
	require.NoError(err)

	event2 := &model.PageViewDuration{
		BID:        "test_update_bid",
		DurationMs: 100,
	}

	err = client.UpdatePageView(ctx, event2)
	assert.NoError(err)
}
