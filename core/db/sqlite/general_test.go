package sqlite_test

import (
	"testing"
)

func TestGetDatabaseVersion(t *testing.T) {
	assert, ctx, client := SetupDatabaseWithUsers(t)

	version, err := client.GetDatabaseVersion(ctx)
	assert.NoError(err)

	assert.Equal("v3.46.0", version)
}
