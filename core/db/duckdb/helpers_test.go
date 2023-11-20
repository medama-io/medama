package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/assert"
)

func SetupDatabase(t *testing.T) (*assert.Assertions, context.Context, *duckdb.Client) {
	t.Helper()
	assert := assert.New(t)
	ctx := context.Background()
	// Disable logging
	log.SetOutput(io.Discard)

	// Generate new memory db per test
	client, err := sqlite.NewClient(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	assert.NoError(err)
	assert.NotNil(client)

	// In memory duckdb client
	duckdbClient, err := duckdb.NewClient("")
	assert.NoError(err)
	assert.NotNil(duckdbClient)

	// Run migrations
	m := migrations.NewMigrationsService(ctx, client, duckdbClient)
	err = m.AutoMigrate(ctx)
	assert.NoError(err)

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
	assert.NoError(err)

	// Create test website
	websiteCreate := model.NewWebsite(
		"duckdb",         // userID
		"medama-test.io", // hostname
		"medama-test.io", // name
		1,                // dateCreated
		2,                // dateUpdated
	)
	err = client.CreateWebsite(ctx, websiteCreate)
	assert.NoError(err)

	return assert, ctx, duckdbClient
}
