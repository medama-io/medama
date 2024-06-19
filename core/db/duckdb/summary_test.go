package duckdb_test

import (
	"testing"

	"github.com/medama-io/medama/db"
)

func TestGetWebsiteSummary(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)
	println("I RAN")
	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    "medium.example.com",
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	println("I RAN2")
	require.NoError(err)

	assert.Equal(871, summary.Visitors)
	assert.Equal(1707, summary.Pageviews)
	assert.Equal(349, summary.Bounces)
	assert.Equal(4994, summary.Duration)
}

/* func testGetWebsiteSummaryEmpty(t *duckdbTest) {
	summary, err := t.client.GetWebsiteSummary(t.ctx, &db.Filters{
		Hostname:    "do-not-exist.example.com",
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	t.require.NoError(err)

	t.assert.Equal(0, summary.Visitors)
	t.assert.Equal(0, summary.Pageviews)
	t.assert.Equal(0, summary.Bounces)
	t.assert.Equal(0, summary.Duration)
}

func testGetWebsiteSummaryFilterAll(t *duckdbTest) {
	summary, err := t.client.GetWebsiteSummary(t.ctx, generateFilterAll())
	t.require.NoError(err)

	t.assert.Equal(871, summary.Visitors)
	t.assert.Equal(1707, summary.Pageviews)
	t.assert.Equal(349, summary.Bounces)
	t.assert.Equal(4994, summary.Duration)
} */
