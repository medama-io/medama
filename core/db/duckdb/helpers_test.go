package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
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
	SimpleFixture Fixture = "./testdata/fixtures/simple.test.db"

	SmallHostname        = "small.example.com"
	MediumHostname       = "medium.example.com"
	DoesNotExistHostname = "does-not-exist.example.com"
)

var (
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TimeStart = time.Unix(0, 0).Format(model.DateFormat)
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TimeEnd = time.Now().Add(24 * time.Hour).Format(model.DateFormat)
)

type SnapRecords struct {
	Records []interface{}
}

type TestCase struct {
	Name    string
	Filters *db.Filters
}

func TestMain(m *testing.M) {
	code := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	// Exit with the same code as the test runner
	os.Exit(code)
}

func NewSnapRecords(slice interface{}) SnapRecords {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Errorf("expected a slice, got %T", slice))
	}

	interfaceSlice := make([]interface{}, v.Len())
	for i := range v.Len() {
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
		"duckdb",                   // userID
		"duckdb@example.com",       // email
		"testtest",                 // password
		model.NewDefaultSettings(), // settings
		1,                          // dateCreated
		2,                          // dateUpdated
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
func generateFilterAll(hostname string) []TestCase {
	filters := make([]TestCase, 0)

	baseFilter := &db.Filters{
		Hostname:    hostname,
		PeriodStart: TimeStart,
		PeriodEnd:   TimeEnd,
	}

	filterSteps := []struct {
		fieldName string
		filter    *db.Filter
	}{
		// Alphabetically sorted list of filters as snapshots are sorted. This lets us see the downward progression of
		// each cumulative filter.
		{"Browser", db.NewFilter(db.FilterBrowser, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("Chrome")}))},
		{"Country", db.NewFilter(db.FilterCountry, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("United Kingdom")}))},
		{"Device", db.NewFilter(db.FilterDevice, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("Desktop")}))},
		{"Language", db.NewFilter(db.FilterLanguage, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("English")}))},
		{"OS", db.NewFilter(db.FilterOS, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("Windows")}))},
		{"Pathname", db.NewFilter(db.FilterPathname, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("/")}))},
		{"Referrer", db.NewFilter(db.FilterReferrer, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("medama.io")}))},
		{"UTMCampaign", db.NewFilter(db.FilterUTMCampaign, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("summer")}))},
		{"UTMMedium", db.NewFilter(db.FilterUTMMedium, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("organic")}))},
		{"UTMSource", db.NewFilter(db.FilterUTMSource, api.NewOptFilterString(api.FilterString{Eq: api.NewOptString("bing")}))},
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
		filters = append(filters, TestCase{Name: step.fieldName, Filters: &tempFilter})
	}

	return filters
}

func getBaseTestCases(_hostname string) []TestCase {
	hostname := MediumHostname // For now we only have one hostname.
	tc := []TestCase{
		{
			Name: "Base",
			Filters: &db.Filters{
				Hostname:    hostname,
				PeriodStart: TimeStart,
				PeriodEnd:   TimeEnd,
			},
		},
		{
			Name: "Empty",
			Filters: &db.Filters{
				Hostname:    DoesNotExistHostname,
				PeriodStart: TimeStart,
				PeriodEnd:   TimeEnd,
			},
		},
	}

	tc = append(tc, generateFilterAll(hostname)...)

	return tc
}
