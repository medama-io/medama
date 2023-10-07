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
	"github.com/stretchr/testify/assert"
)

func SetupDatabase(t *testing.T) (*assert.Assertions, context.Context, *sqlite.Client) {
	t.Helper()
	assert := assert.New(t)
	ctx := context.Background()
	// Disable logging
	log.SetOutput(io.Discard)

	client, err := sqlite.NewClient(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()))
	assert.NoError(err)
	assert.NotNil(client)

	// Run migrations
	m := migrations.NewMigrationsService(client)
	err = m.AutoMigrate()
	assert.NoError(err)

	return assert, ctx, client
}
