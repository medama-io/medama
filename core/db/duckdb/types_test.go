package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsiteBrowsersSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteBrowsersSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteBrowsers(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			browsers, err := client.GetWebsiteBrowsers(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(browsers)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteOSSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteOSSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteOS(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			os, err := client.GetWebsiteOS(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(os)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteDevicesSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteDevicesSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteDevices(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			devices, err := client.GetWebsiteDevices(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(devices)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
