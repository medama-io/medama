package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/model"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/vfs/memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupDatabase(t *testing.T) (*assert.Assertions, context.Context, *duckdb.Client) {
	t.Helper()
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()
	// Disable logging
	log.SetOutput(io.Discard)

	// Generate new memory db per test
	name := fmt.Sprintf("file:/%s.db?vfs=memdb", t.Name())
	memdb.Create(name, []byte{})
	client, err := sqlite.NewClient(name)
	require.NoError(err)
	assert.NotNil(client)

	// In memory duckdb client
	duckdbClient, err := duckdb.NewClient("")
	require.NoError(err)
	assert.NotNil(duckdbClient)

	// Run migrations
	m := migrations.NewMigrationsService(ctx, client, duckdbClient)
	err = m.AutoMigrate(ctx)
	require.NoError(err)

	// Create test user
	userCreate := model.NewUser(
		"duckdb",             // userID
		"duckdb@example.com", // email
		"testtest",           // password
		"en",                 // language
		1,                    // dateCreated
		2,                    // dateUpdated
	)
	err = client.CreateUser(ctx, userCreate)
	require.NoError(err)

	// Create test website
	websiteCreate := model.NewWebsite(
		"duckdb",         // userID
		"medama-test.io", // hostname
		"medama-test.io", // name
		1,                // dateCreated
		2,                // dateUpdated
	)
	err = client.CreateWebsite(ctx, websiteCreate)
	require.NoError(err)

	return assert, ctx, duckdbClient
}
