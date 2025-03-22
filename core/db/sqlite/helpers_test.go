package sqlite_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	_ "github.com/marcboeker/go-duckdb/v2"
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

func SetupDatabase(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
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

	// Empty duckdb client not used in tests
	duckdbClient, err := duckdb.NewClient("")
	require.NoError(err)
	assert.NotNil(duckdbClient)

	// Run migrations
	m, err := migrations.NewMigrationsService(ctx, client, duckdbClient)
	require.NoError(err)
	err = m.AutoMigrate(ctx)
	assert.NoError(err)

	return assert, ctx, client
}

func SetupDatabaseWithUsers(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
	t.Helper()
	assert, ctx, client := SetupDatabase(t)

	ids := []string{"test1", "test2", "test3"}
	usernames := []string{"username1", "username2", "username3"}
	passwords := []string{"password1", "password2", "password3"}

	for i, id := range ids {
		userCreate := model.NewUser(
			id,
			usernames[i],
			passwords[i],
			model.NewDefaultSettings(),
			1,
			2,
		)

		err := client.CreateUser(ctx, userCreate)
		assert.NoError(err)
	}

	return assert, ctx, client
}

func SetupDatabaseWithWebsites(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
	t.Helper()
	assert, ctx, client := SetupDatabaseWithUsers(t)

	ids := []string{"website1", "website2", "website3"}
	userIDs := []string{"test1", "test2", "test3"}

	// 3 websites each for 3 users
	for _, id := range ids {
		for _, userID := range userIDs {
			websiteCreate := model.NewWebsite(
				userID,
				fmt.Sprintf("%s-%s.com", id, userID),
				1,
				2,
			)

			err := client.CreateWebsite(ctx, websiteCreate)
			assert.NoError(err)
		}
	}

	return assert, ctx, client
}
