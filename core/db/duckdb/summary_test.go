package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/medama-io/medama/db"
)

func TestGetWebsiteSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)
	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    MediumHostname,
		PeriodStart: TimeStart,
		PeriodEnd:   TimeEnd,
	})
	require.NoError(err)

	snaps.MatchJSON(t, summary)
}

func TestGetWebsiteSummaryEmpty(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    DoesNotExistHostname,
		PeriodStart: TimeStart,
		PeriodEnd:   TimeEnd,
	})
	require.NoError(err)

	assert.Equal(0, summary.Visitors)
	assert.Equal(0, summary.Pageviews)
	assert.Equal(float32(0), summary.BounceRate)
	assert.Equal(0, summary.Duration)
}

func TestGetWebsiteSummaryFilterAll(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	for _, filter := range generateFilterAll(MediumHostname) {
		summary, err := client.GetWebsiteSummary(ctx, filter.Filters)
		require.NoError(err)

		snaps.MatchJSON(t, summary)
	}
}

// TODO: Need to make a new deterministic time fixture for these tests.
// func TestGetWebsiteIntervals(t *testing.T) {}
// func TestGetWebsiteLast24Hours(t *testing.T) {}
