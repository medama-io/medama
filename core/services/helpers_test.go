package services_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/services"
	"github.com/medama-io/medama/util"
	_ "github.com/ncruces/go-sqlite3/driver"
	"github.com/ncruces/go-sqlite3/vfs/memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestHandler(
	t *testing.T,
) (*assert.Assertions, context.Context, *services.Handler, *sqlite.Client) {
	t.Helper()
	return newTestHandler(t, false)
}

func NewTestHandlerDemoMode(
	t *testing.T,
) (*assert.Assertions, context.Context, *services.Handler, *sqlite.Client) {
	t.Helper()
	return newTestHandler(t, true)
}

func newTestHandler(
	t *testing.T,
	isDemoMode bool,
) (*assert.Assertions, context.Context, *services.Handler, *sqlite.Client) {
	t.Helper()

	assert := assert.New(t)
	require := require.New(t)
	ctx := t.Context()

	log.SetOutput(io.Discard)

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

	auth, err := util.NewAuthService(ctx, isDemoMode)
	require.NoError(err)

	handler, err := services.NewService(ctx, auth, sqliteClient, duckdbClient, "test-commit")
	require.NoError(err)
	assert.NotNil(handler)

	return assert, ctx, handler, sqliteClient
}
