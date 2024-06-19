package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/medama-io/medama/db"
)

func TestGetWebsitePagesSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)
	summary, err := client.GetWebsitePagesSummary(ctx, &db.Filters{
		Hostname:    MEDIUM_HOSTNAME,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	snap := NewSnapRecords(summary)
	snaps.MatchSnapshot(t, snap.Snapshot())
}

func TestGetWebsitePagesSummaryEmpty(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	summary, err := client.GetWebsitePagesSummary(ctx, &db.Filters{
		Hostname:    DOES_NOT_EXIST_HOSTNAME,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	assert.Empty(summary)
}

func TestGetWebsitePagesSummaryFilterAll(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	for _, filter := range generateFilterAll(MEDIUM_HOSTNAME) {
		summary, err := client.GetWebsitePagesSummary(ctx, filter)
		require.NoError(err)

		snap := NewSnapRecords(summary)
		snaps.MatchSnapshot(t, snap.Snapshot())
	}
}

func TestGetWebsitePages(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)
	pages, err := client.GetWebsitePages(ctx, &db.Filters{
		Hostname:    MEDIUM_HOSTNAME,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	snap := NewSnapRecords(pages)
	snaps.MatchSnapshot(t, snap.Snapshot())
}

func TestGetWebsitePagesEmpty(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	pages, err := client.GetWebsitePages(ctx, &db.Filters{
		Hostname:    DOES_NOT_EXIST_HOSTNAME,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	assert.Empty(pages)
}

func TestGetWebsitePagesFilterAll(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	for _, filter := range generateFilterAll(MEDIUM_HOSTNAME) {
		pages, err := client.GetWebsitePages(ctx, filter)
		require.NoError(err)

		snap := NewSnapRecords(pages)
		snaps.MatchSnapshot(t, snap.Snapshot())
	}
}
