package metest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/ncruces/go-sqlite3/vfs/memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewInMemoryDatabase(t *testing.T) (*sqlite.Client, *duckdb.Client) {
	t.Helper()

	assert := assert.New(t)
	require := require.New(t)
	ctx := t.Context()

	name := strings.ToLower(t.Name() + "_svc_test")
	host := fmt.Sprintf("file:/%s.db?vfs=memdb", name)

	memdb.Create(name, []byte{})

	sqliteClient, err := sqlite.NewClient(host)
	require.NoError(err)
	assert.NotNil(sqliteClient)

	duckdbClient, err := duckdb.NewClient(":memory:")
	require.NoError(err)
	assert.NotNil(duckdbClient)

	m, err := migrations.NewMigrationsService(ctx, sqliteClient, duckdbClient)
	require.NoError(err)
	err = m.AutoMigrate(ctx)
	require.NoError(err)

	return sqliteClient, duckdbClient
}
