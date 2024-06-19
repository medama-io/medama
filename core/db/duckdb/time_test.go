package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsiteTimeSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteTimeSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteTime(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			times, err := client.GetWebsiteTime(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(times)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
