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
	t.assert.NoError(err)

	t.assert.Equal(871, summary.Visitors)
	t.assert.Equal(1707, summary.Pageviews)
	t.assert.Equal(349, summary.Bounces)
	t.assert.Equal(4994, summary.Duration)
}
