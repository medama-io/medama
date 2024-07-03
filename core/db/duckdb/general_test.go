package duckdb_test

import (
	"testing"
)

func TestGetSettingsUsage(t *testing.T) {
	assert, require, ctx, client := UseDatabaseFixture(t, SIMPLE_FIXTURE)

	usage, err := client.GetSettingsUsage(ctx)
	require.NoError(err)

	assert.NotNil(usage.Threads)
	assert.NotNil(usage.MemoryLimit)
}
