package duckdb_test

import (
	"github.com/medama-io/medama/db"
)

func testGetWebsiteSummary(t *duckdbTest) {
	summary, err := t.client.GetWebsiteSummary(t.ctx, &db.Filters{
		Hostname:    "1.example.com",
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	})
	t.require.NoError(err)

	t.assert.Equal(871, summary.Visitors)
	t.assert.Equal(1707, summary.Pageviews)
	t.assert.Equal(349, summary.Bounces)
	t.assert.Equal(4994, summary.Duration)
}

func testGetWebsiteSummaryEmpty(t *duckdbTest) {
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
}
