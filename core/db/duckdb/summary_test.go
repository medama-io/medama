package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/medama-io/medama/db"
)

func TestGetWebsiteSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)
	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    MEDIUM_HOSTNAME,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	snaps.MatchJSON(t, summary)
}

func TestGetWebsiteSummaryEmpty(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    "do-not-exist.example.com",
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	require.NoError(err)

	assert.Equal(0, summary.Visitors)
	assert.Equal(0, summary.Pageviews)
	assert.Equal(0, summary.Bounces)
	assert.Equal(0, summary.Duration)
}

func TestGetWebsiteSummaryFilterAll(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	summary, err := client.GetWebsiteSummary(ctx, generateFilterAll(MEDIUM_HOSTNAME))
	require.NoError(err)

	snaps.MatchJSON(t, summary)
}
