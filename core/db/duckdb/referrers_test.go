package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsiteReferrersSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteReferrersSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteReferrers(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			referrers, err := client.GetWebsiteReferrers(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(referrers)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
