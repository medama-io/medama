package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
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

	SMALL_HOSTNAME          = "small.example.com"
	MEDIUM_HOSTNAME         = "medium.example.com"
	DOES_NOT_EXIST_HOSTNAME = "does-not-exist.example.com"
)

var (
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_START = time.Unix(0, 0).Format(model.DateFormat)
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_END = time.Now().Add(24 * time.Hour).Format(model.DateFormat)
)

type SnapRecords struct {
	Records []interface{}
}

func NewSnapRecords(slice interface{}) SnapRecords {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Errorf("expected a slice, got %T", slice))
	}

	interfaceSlice := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		interfaceSlice[i] = v.Index(i).Interface()
	}

	return SnapRecords{Records: interfaceSlice}
}

// Generate a new snapshot for the given records.
func (s SnapRecords) Snapshot() string {
	var sb strings.Builder
	sb.WriteString("RECORDS:\n")
	for _, record := range s.Records {
		sb.WriteString(fmt.Sprintf("%+v\n", record))
	}
	return sb.String()
}

func TestMain(m *testing.M) {
	m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
}

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

// Generate an array of filters where we incrementally add one new filter to the previous one.
func generateFilterAll(hostname string) []*db.Filters {
	filters := make([]*db.Filters, 0)

	baseFilter := &db.Filters{
		Hostname:    hostname,
		PeriodStart: TIME_START,
		PeriodEnd:   TIME_END,
	}
	// Make a copy to avoid receiving future modifications.
	tempFilter := *baseFilter
	filters = append(filters, &tempFilter)

	filterSteps := []struct {
		fieldName string
		filter    *db.Filter
	}{
		// Inverted list of least importance to most importance for most appropriate and relevant snapshot results.
		{"UTMSource", db.NewFilter(db.FilterUTMSource, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("bing")}))},
		{"UTMMedium", db.NewFilter(db.FilterUTMMedium, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("organic")}))},
		{"UTMCampaign", db.NewFilter(db.FilterUTMCampaign, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("summer")}))},

		{"Language", db.NewFilter(db.FilterLanguage, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("en")}))},
		{"Country", db.NewFilter(db.FilterCountry, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("GB")}))},

		{"Device", db.NewFilter(db.FilterDevice, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Desktop")}))},
		{"OS", db.NewFilter(db.FilterOS, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Windows")}))},
		{"Browser", db.NewFilter(db.FilterBrowser, api.NewOptFilterFixed(api.FilterFixed{Eq: api.NewOptString("Chrome")}))},
		{"Referrer", db.NewFilter(db.FilterReferrer, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("medama.io")}))},
		{"Pathname", db.NewFilter(db.FilterPathname, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("/")}))},
	}

	for _, step := range filterSteps {
		// Modify the base filter with the new filter.
		switch step.fieldName {
		case "Pathname":
			baseFilter.Pathname = step.filter
		case "Referrer":
			baseFilter.Referrer = step.filter
		case "UTMSource":
			baseFilter.UTMSource = step.filter
		case "UTMMedium":
			baseFilter.UTMMedium = step.filter
		case "UTMCampaign":
			baseFilter.UTMCampaign = step.filter
		case "Browser":
			baseFilter.Browser = step.filter
		case "OS":
			baseFilter.OS = step.filter
		case "Device":
			baseFilter.Device = step.filter
		case "Country":
			baseFilter.Country = step.filter
		case "Language":
			baseFilter.Language = step.filter
		}

		// Make a local copy to avoid overwriting the previous filter.
		tempFilter := *baseFilter
		filters = append(filters, &tempFilter)
	}

	return filters
}
