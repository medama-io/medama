package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsiteReferrersSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	testCases := getBaseTestCases(MediumHostname)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteReferrersSummary(ctx, false, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteReferrersSummaryGrouped(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	testCases := getBaseTestCases(MediumHostname)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteReferrersSummary(ctx, true, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteReferrers(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	testCases := getBaseTestCases(MediumHostname)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			referrers, err := client.GetWebsiteReferrers(ctx, false, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(referrers)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteReferrersGrouped(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SimpleFixture)

	testCases := getBaseTestCases(MediumHostname)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			referrers, err := client.GetWebsiteReferrers(ctx, true, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(referrers)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
