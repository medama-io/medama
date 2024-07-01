package duckdb_test

import (
	"testing"
)

func TestGetDatabaseVersion(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	version, err := client.GetDatabaseVersion(ctx)
	require.NoError(err)

	assert.Equal("v1.0.0", version)
}
