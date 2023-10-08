package sqlite_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
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

	// Run migrations
	m := migrations.NewMigrationsService(ctx, client)
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
		userCreate := &model.User{
			ID:          id,
			Email:       emails[i],
			Password:    passwords[i],
			Language:    "en",
			DateCreated: 1,
			DateUpdated: 2,
		}

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
			websiteCreate := &model.Website{
				ID:          fmt.Sprintf("%s-%s", id, user_id),
				UserID:      user_id,
				Hostname:    fmt.Sprintf("%s-%s.com", id, user_id),
				DateCreated: 1,
				DateUpdated: 2,
			}

			err := client.CreateWebsite(ctx, websiteCreate)
			assert.NoError(err)
		}
	}

	return assert, ctx, client
}
