package main

import (
	"context"
	"math/rand/v2"
	"time"

	_ "github.com/marcboeker/go-duckdb"
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

	// Generate fixtures.
	//nolint: gosec // Using a fixed seed will produce the same output on every run.
	smallRand := rand.New(rand.NewPCG(1, 2))
	generateFixture(smallRand, SMALL_FIXTURE_COUNT, hostnames[0], duckdb)
	//nolint: gosec // Using a fixed seed will produce the same output on every run.
	mediumRand := rand.New(rand.NewPCG(1, 2))
	generateFixture(mediumRand, MEDIUM_FIXTURE_COUNT, hostnames[1], duckdb)
}

func generateFixture(r *rand.Rand, count int, hostname string, client *duckdb.Client) {
	// We should use the Appender API for this, but we were encountering corruption issues and segfaults
	// from the go-duckdb driver. We can revisit this in the future once it is more stable.
	//
	// This is an order of magnitude slower than the Appender API, but it is good enough for now.
	pageViewHits := generatePageViewHits(r, count, hostname)
	for _, pv := range pageViewHits {
		err := client.AddPageView(context.Background(), pv)
		if err != nil {
			panic(err)
		}

	}

	pageViewDurations := generatePageViewDurations(r, pageViewHits, count)
	for _, pd := range pageViewDurations {
		err := client.AddPageDuration(context.Background(), pd)
		if err != nil {
			panic(err)
		}
	}
}
