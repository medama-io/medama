package duckdb_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"testing"
	"time"

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

const (
	// Number of page views to generate.
	PAGEVIEW_COUNT = 5000
	// Number of page view durations to generate.
	// Account for unreliability of failing to send page durations with a lower value.
	DURATION_COUNT = 4000
)

var (
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_START = time.Unix(0, 0).Format(model.DateFormat)
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_END = time.Now().Add(24 * time.Hour).Format(model.DateFormat)
)

func SetupDatabase(t *testing.T) *duckdbTest {
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

	// Create test website
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

	return &duckdbTest{assert, require, ctx, duckdbClient}
}

func SetupDatabaseWithPageViews(t *testing.T) *duckdbTest {
	t.Helper()
	db := SetupDatabase(t)

	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewPCG(1, 2))

	// Generate page view hits.
	pageViewHits := generatePageViewHits(r, PAGEVIEW_COUNT)
	for _, pageViewHit := range pageViewHits {
		err := db.client.AddPageView(db.ctx, pageViewHit)
		db.require.NoError(err)
	}

	// Generate page view durations.
	pageViewDurations := generatePageViewDurations(r, pageViewHits, DURATION_COUNT)
	for _, pageViewDuration := range pageViewDurations {
		err := db.client.AddPageDuration(db.ctx, pageViewDuration)
		db.require.NoError(err)
	}

	return db
}

func generatePageViewHits(r *rand.Rand, count int) []*model.PageViewHit {
	hostnames := []string{"1.example.com", "2.example.com", "3.example.com"}
	paths := []string{"/", "/about", "/contact", "/pricing", "/blog"}
	booleanValues := []bool{true, false}
	referrers := []string{"1.example.com", "medama.io", "google.com"}
	countryCodes := []string{"GB", "US", "DE", "FR", "ES", "IT"}
	languages := []string{"en", "de", "fr", "es", "it"}
	browserNames := []model.BrowserName{model.ChromeBrowser, model.FirefoxBrowser, model.SafariBrowser, model.EdgeBrowser}
	oses := []model.OSName{model.WindowsOS, model.MacOS, model.LinuxOS, model.AndroidOS, model.IOS}
	deviceTypes := []model.DeviceType{model.DesktopDevice, model.MobileDevice, model.TabletDevice}
	utmSources := []string{"", "bing", "twitter"}
	utmMediums := []string{"", "cpc", "organic"}
	utmCampaigns := []string{"", "summer", "winter"}

	pageViewHits := make([]*model.PageViewHit, count)

	for i := range count {
		pageViewHits[i] = &model.PageViewHit{
			Hostname:     hostnames[r.IntN(len(hostnames))],
			Pathname:     paths[r.IntN(len(paths))],
			IsUniqueUser: booleanValues[r.IntN(len(booleanValues))],
			IsUniquePage: booleanValues[r.IntN(len(booleanValues))],
			Referrer:     referrers[r.IntN(len(referrers))],
			CountryCode:  countryCodes[r.IntN(len(countryCodes))],
			Language:     languages[r.IntN(len(languages))],
			BrowserName:  browserNames[r.IntN(len(browserNames))],
			OS:           oses[r.IntN(len(oses))],
			DeviceType:   deviceTypes[r.IntN(len(deviceTypes))],
			UTMSource:    utmSources[r.IntN(len(utmSources))],
			UTMMedium:    utmMediums[r.IntN(len(utmMediums))],
			UTMCampaign:  utmCampaigns[r.IntN(len(utmCampaigns))],
		}
	}

	return pageViewHits
}

func generatePageViewDurations(r *rand.Rand, hits []*model.PageViewHit, count int) []*model.PageViewDuration {
	durations := make([]*model.PageViewDuration, count)

	for i := range count {
		durations[i] = &model.PageViewDuration{
			PageViewHit: *hits[i],
			DurationMs:  r.IntN(10000),
		}
	}

	return durations
}
