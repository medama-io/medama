package duckdb_test

import (
	"testing"
	"time"

	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

func TestGetWebsiteSummary(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithPageViews(t)
	defer client.Close()

	summary, err := client.GetWebsiteSummary(ctx, &db.Filters{
		Hostname:    "1.example.com",
		PeriodStart: time.Unix(0, 0).Format(model.DateFormat),
		PeriodEnd:   time.Now().Format(model.DateFormat),
	})
	assert.NoError(err)

	assert.Equal(10, summary.Visitors)
	assert.Equal(10, summary.Pageviews)
	assert.Equal(10, summary.Bounces)
	assert.Equal(10, summary.Duration)
}
