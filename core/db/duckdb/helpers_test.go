package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
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

type Fixture string

const (
	SIMPLE_FIXTURE Fixture = "./testdata/fixtures/simple.test.db"
)

var (
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_START = time.Unix(0, 0).Format(model.DateFormat)
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_END = time.Now().Add(24 * time.Hour).Format(model.DateFormat)
)

func SetupDatabase(t *testing.T) (*assert.Assertions, *require.Assertions, context.Context, *duckdb.Client) {
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
	require.NoError(client.Ping())
	assert.NotNil(client)

	// In memory duckdb client.
	duckdbClient, err := duckdb.NewClient("")
	require.NoError(err)
	require.NoError(duckdbClient.Ping())
	assert.NotNil(duckdbClient)

	// Run migrations
	m, err := migrations.NewMigrationsService(ctx, client, duckdbClient)
	require.NoError(err)
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

	// Create test websites.
	hostnames := []string{"1.example.com", "2.example.com", "3.example.com"}
	for _, hostname := range hostnames {
		websiteCreate := model.NewWebsite(
			"duckdb", // userID
			hostname, // hostname
			1,        // dateCreated
			2,        // dateUpdated
		)
		err = client.CreateWebsite(ctx, websiteCreate)
		require.NoError(err)
	}

	return assert, require, ctx, duckdbClient
}

func UseDatabaseFixture(t *testing.T, fixture Fixture) (*assert.Assertions, *require.Assertions, context.Context, *duckdb.Client) {
	t.Helper()
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()

	client, err := duckdb.NewClient(string(fixture))
	require.NoError(err)
	require.NoError(client.Ping())
	assert.NotNil(client)

	return assert, require, ctx, client
}

func generateFilterAll(hostname string) *db.Filters {
	return &db.Filters{
		Hostname:    hostname,
		Pathname:    db.NewFilter(db.FilterPathname, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("/")})),
		Referrer:    db.NewFilter(db.FilterReferrer, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("medama.io")})),
		UTMSource:   db.NewFilter(db.FilterUTMSource, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("bing")})),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("organic")})),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("summer")})),
		// Browser:     db.NewFilter(db.FilterBrowser, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Chrome")})),
		// OS:          db.NewFilter(db.FilterOS, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Windows")})),
		// Device:      db.NewFilter(db.FilterDevice, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Desktop")})),
		// Country:     db.NewFilter(db.FilterCountry, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("GB")})),
		// Language:    db.NewFilter(db.FilterLanguage, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("en")})),

		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	}
}
