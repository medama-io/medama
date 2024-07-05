package duckdb_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestGetWebsiteCountriesSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteCountriesSummary(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteCountries(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			countries, err := client.GetWebsiteCountries(ctx, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(countries)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteLanguagesSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteLanguagesSummary(ctx, false, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteLanguagesLocaleSummary(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			summary, err := client.GetWebsiteLanguagesSummary(ctx, true, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(summary)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteLanguages(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			languages, err := client.GetWebsiteLanguages(ctx, false, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(languages)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}

func TestGetWebsiteLanguagesLocale(t *testing.T) {
	_, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	testCases := getBaseTestCases(MEDIUM_HOSTNAME)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			languages, err := client.GetWebsiteLanguages(ctx, true, tc.Filters)
			require.NoError(err)

			snap := NewSnapRecords(languages)
			snaps.MatchSnapshot(t, snap.Snapshot())
		})
	}
}
