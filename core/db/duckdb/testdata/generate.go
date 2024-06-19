package main

import (
	"context"
	"database/sql/driver"
	"math/rand/v2"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	goduckdb "github.com/marcboeker/go-duckdb"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/model"
	"github.com/ncruces/go-sqlite3/vfs/memdb"
)

const (
	SMALL_FIXTURE_COUNT  = 10000
	MEDIUM_FIXTURE_COUNT = 1000000

	sqliteHost = "file:/fixture.db?vfs=memdb"
	duckdbHost = "./db/duckdb/testdata/fixtures/simple.test.db"
)

var (
	// Generate a 1 month interval of data between 1st Jan 2024 and 1st Feb 2024.
	//
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_START = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	//nolint:gochecknoglobals // Reason: These are used in every test.
	TIME_END = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
)

func main() {
	ctx := context.Background()

	// SQLite client to run migrations.
	memdb.Create("temp", []byte{})
	sqlite, err := sqlite.NewClient(sqliteHost)
	if err != nil {
		panic(err)
	}
	defer sqlite.Close()
	if err := sqlite.Ping(); err != nil {
		panic(err)
	}

	// DuckDB client.
	duckdb, err := duckdb.NewClient(duckdbHost)
	if err != nil {
		panic(err)
	}
	defer duckdb.Close()
	if err := duckdb.Ping(); err != nil {
		panic(err)
	}

	// Run migrations.
	m, err := migrations.NewMigrationsService(ctx, sqlite, duckdb)
	if err != nil {
		panic(err)
	}
	if err := m.AutoMigrate(ctx); err != nil {
		panic(err)
	}

	// Create test user.
	userCreate := model.NewUser(
		"duckdb",             // userID
		"duckdb@example.com", // email
		"testtest",           // password
		"en",                 // language
		1,                    // dateCreated
		2,                    // dateUpdated
	)

	if err := sqlite.CreateUser(ctx, userCreate); err != nil {
		panic(err)
	}

	// Create test websites.
	hostnames := []string{"small.example.com", "medium.example.com"}
	for _, hostname := range hostnames {
		websiteCreate := model.NewWebsite(
			"duckdb", // userID
			hostname, // hostname
			1,        // dateCreated
			2,        // dateUpdated
		)
		err = sqlite.CreateWebsite(ctx, websiteCreate)
		if err != nil {
			panic(err)
		}
	}

	// Initialise appender API for quicker inserts.
	connector, err := goduckdb.NewConnector(duckdbHost, nil)
	if err != nil {
		panic(err)
	}

	conn, err := connector.Connect(ctx)
	if err != nil {
		panic(err)
	}

	viewAppender, err := goduckdb.NewAppenderFromConn(conn, "", "views")
	if err != nil {
		panic(err)
	}

	durationAppender, err := goduckdb.NewAppenderFromConn(conn, "", "duration")
	if err != nil {
		panic(err)
	}

	// Generate fixtures.
	//nolint: gosec // Using a fixed seed will produce the same output on every run.
	smallRand := rand.New(rand.NewPCG(1, 2))
	generateFixture(smallRand, SMALL_FIXTURE_COUNT, hostnames[0], viewAppender, durationAppender)
	//nolint: gosec // Using a fixed seed will produce the same output on every run.
	mediumRand := rand.New(rand.NewPCG(1, 2))
	generateFixture(mediumRand, MEDIUM_FIXTURE_COUNT, hostnames[1], viewAppender, durationAppender)

	// Close appenders
	err = viewAppender.Close()
	if err != nil {
		panic(err)
	}

	err = durationAppender.Close()
	if err != nil {
		panic(err)
	}
	err = conn.Close()
	if err != nil {
		panic(err)
	}
	err = connector.Close()
	if err != nil {
		panic(err)
	}
}

func generateFixture(r *rand.Rand, count int, hostname string, viewAppender *goduckdb.Appender, durationAppender *goduckdb.Appender) {
	pageViewHits := generatePageViewHits(r, count, hostname)
	for _, pv := range pageViewHits {
		if err := viewAppender.AppendRow(pageViewToValues(r, pv)...); err != nil {
			panic(err)
		}
	}

	err := viewAppender.Flush()
	if err != nil {
		panic(err)
	}

	pageViewDurations := generatePageViewDurations(r, pageViewHits, count)
	for _, pd := range pageViewDurations {
		if err := durationAppender.AppendRow(pageDurationToValues(r, pd)...); err != nil {
			panic(err)
		}
	}

	err = durationAppender.Flush()
	if err != nil {
		panic(err)
	}
}

func pageViewToValues(r *rand.Rand, pv *model.PageViewHit) []driver.Value {
	return []driver.Value{
		pv.Hostname,
		pv.Pathname,
		pv.IsUniqueUser,
		pv.IsUniquePage,
		pv.Referrer,
		pv.CountryCode,
		pv.Language,
		uint8(pv.BrowserName),
		uint8(pv.OS),
		uint8(pv.DeviceType),
		pv.UTMSource,
		pv.UTMMedium,
		pv.UTMCampaign,
		// Generate a random time.Time between TIME_START and TIME_END.
		TIME_START.Add(time.Duration(r.Int64N(TIME_END.Unix() - TIME_START.Unix()))),
	}
}

func pageDurationToValues(r *rand.Rand, pd *model.PageViewDuration) []driver.Value {
	return []driver.Value{
		pd.Hostname,
		pd.Pathname,
		pd.IsUniqueUser,
		pd.IsUniquePage,
		pd.Referrer,
		pd.CountryCode,
		pd.Language,
		uint8(pd.BrowserName),
		uint8(pd.OS),
		uint8(pd.DeviceType),
		pd.UTMSource,
		pd.UTMMedium,
		pd.UTMCampaign,
		int32(pd.DurationMs),
		// Generate a random time.Time between TIME_START and TIME_END.
		TIME_START.Add(time.Duration(r.Int64N(TIME_END.Unix() - TIME_START.Unix()))),
	}
}
