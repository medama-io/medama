package duckdb_test

import (
	"context"
	"testing"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type duckdbTest struct {
	assert  *assert.Assertions
	require *require.Assertions
	//nolint: containedctx // This is a test filler context with no impact.
	ctx    context.Context
	client *duckdb.Client
}

func TestAll(t *testing.T) {
	// go-duckdb does not play well with many in-memory or file based tests, so we
	// need to run the tests with a single client. Almost all tests are read-only,
	// so we can reuse the same client for all tests.
	//
	// Hopefully, this will become more stable in the future and we can go back to
	// more idiomatic testing.
	db := SetupDatabaseWithPageViews(t)
	defer db.client.Close()

	tests := map[string]func(*duckdbTest){
		"TestAddPageView":     testAddPageView,
		"TestAddPageDuration": testAddPageDuration,
		// Summary
		"TestGetWebsiteSummary":          testGetWebsiteSummary,
		"TestGetWebsiteSummaryEmpty":     testGetWebsiteSummaryEmpty,
		"TestGetWebsiteSummaryFilterAll": testGetWebsiteSummaryFilterAll,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test(db)
		})
	}
}
