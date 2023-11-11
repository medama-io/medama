package sqlite_test

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

func SetupDatabase(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
	t.Helper()
	assert := assert.New(t)
	ctx := context.Background()
	// Disable logging
	log.SetOutput(io.Discard)

	// Generate new memory db per test
	client, err := sqlite.NewClient(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	assert.NoError(err)
	assert.NotNil(client)

	// Empty duckdb client not used in tests
	duckdbClient, err := duckdb.NewClient("")
	assert.NoError(err)
	assert.NotNil(duckdbClient)

	// Run migrations
	m := migrations.NewMigrationsService(ctx, client, duckdbClient)
	err = m.AutoMigrate(ctx)
	assert.NoError(err)

	return assert, ctx, client
}

func SetupDatabaseWithUsers(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
	t.Helper()
	assert, ctx, client := SetupDatabase(t)

	ids := []string{"test1", "test2", "test3"}
	emails := []string{"test1@example.com", "test2@example.com", "test3@example.com"}
	passwords := []string{"password1", "password2", "password3"}

	for i, id := range ids {
		userCreate := model.NewUser(
			id,
			emails[i],
			passwords[i],
			"en",
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
	user_ids := []string{"test1", "test2", "test3"}

	// 3 websites each for 3 users
	for _, id := range ids {
		for _, user_id := range user_ids {
			websiteCreate := model.NewWebsite(
				user_id,
				fmt.Sprintf("%s-%s.com", id, user_id),
				fmt.Sprintf("%s-%s", id, user_id),
				1,
				2,
			)

			err := client.CreateWebsite(ctx, websiteCreate)
			assert.NoError(err)
		}
	}

	return assert, ctx, client
}
