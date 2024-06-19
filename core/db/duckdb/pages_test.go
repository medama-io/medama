package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsitePagesSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsitePagesSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsitePages(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			pages, err := client.GetWebsitePages(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(pages)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
