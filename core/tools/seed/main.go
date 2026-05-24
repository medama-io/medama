package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"github.com/rs/zerolog"
)

const (
	defaultAppDB       = "./me_meta.db"
	defaultAnalyticsDB = "./me_analytics.db"
	defaultDays        = 60
)

type config struct {
	appDB       string
	analyticsDB string
	size        string
	reset       bool
	days        int
}

type seedOwner struct {
	ID       string `db:"id"`
	Username string `db:"username"`
}

type site struct {
	hostname string
	weight   int
	paths    []string
}

type language struct {
	base    string
	dialect string
}

type property struct {
	name   string
	values []string
}

var sites = []site{
	{
		hostname: "localhost",
		weight:   5,
		paths:    []string{"/", "/pricing", "/docs/install", "/settings/account"},
	},
	{
		hostname: "docs.medama.test",
		weight:   3,
		paths:    []string{"/", "/getting-started", "/reference/api", "/features/custom-properties"},
	},
	{
		hostname: "shop.medama.test",
		weight:   2,
		paths:    []string{"/", "/products", "/products/pro", "/checkout"},
	},
}

var (
	countries        = []string{"United States", "Japan", "United Kingdom", "Germany", "Canada"}
	referrers        = []string{"", "google.com", "github.com", "news.ycombinator.com", "docs.medama.io"}
	languages        = []language{{"English", "American English"}, {"English", "British English"}, {"Japanese", "Japanese"}, {"German", "German"}}
	browsers         = []string{"Chrome", "Safari", "Firefox", "Edge"}
	operatingSystems = []string{"macOS", "Windows", "iOS", "Android", "Linux"}
	devices          = []string{"Desktop", "Mobile", "Tablet"}
	utmSources       = []string{"newsletter", "google", "github", "x", "partner"}
	utmMediums       = []string{"email", "organic", "social", "cpc", "referral"}
	utmCampaigns     = []string{"spring-launch", "product-update", "self-hosted", "privacy-tools"}
	properties       = []property{
		{"plan", []string{"free", "starter", "pro", "team"}},
		{"theme", []string{"system", "dark", "light"}},
		{"cta", []string{"add-website", "copy-snippet", "view-docs", "upgrade"}},
	}
)

func main() {
	cfg := parseFlags()
	if err := run(context.Background(), cfg); err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.appDB, "appdb", defaultAppDB, "Path to the SQLite app database.")
	flag.StringVar(&cfg.analyticsDB, "analyticsdb", defaultAnalyticsDB, "Path to the DuckDB analytics database.")
	flag.StringVar(&cfg.size, "size", "small", "Fixture size: small or medium.")
	flag.BoolVar(&cfg.reset, "reset", false, "Delete existing data for seed hostnames before seeding.")
	flag.IntVar(&cfg.days, "days", defaultDays, "Number of recent days to spread analytics across.")
	flag.Parse()
	return cfg
}

func run(ctx context.Context, cfg config) error {
	if cfg.days <= 0 {
		return fmt.Errorf("days must be greater than 0")
	}

	views := 0
	switch cfg.size {
	case "small":
		views = 10_000
	case "medium":
		views = 250_000
	default:
		return fmt.Errorf("unknown size %q; expected small or medium", cfg.size)
	}

	if _, err := logger.Init("pretty", "info"); err != nil {
		return err
	}
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	appDB, err := sqlite.NewClient(cfg.appDB)
	if err != nil {
		return fmt.Errorf("open app db: %w", err)
	}
	defer appDB.Close()

	analyticsDB, err := duckdb.NewClient(cfg.analyticsDB)
	if err != nil {
		return fmt.Errorf("open analytics db: %w", err)
	}
	defer analyticsDB.Close()

	migrator, err := migrations.NewMigrationsService(ctx, appDB, analyticsDB)
	if err != nil {
		return fmt.Errorf("create migrations service: %w", err)
	}
	if err := migrator.AutoMigrate(ctx); err != nil {
		return fmt.Errorf("migrate databases: %w", err)
	}

	user, err := getSeedOwner(ctx, appDB)
	if err != nil {
		return err
	}

	hostnames := make([]string, len(sites))
	for i, site := range sites {
		hostnames[i] = site.hostname
	}

	exists, err := hasSeedData(ctx, appDB, analyticsDB, hostnames)
	if err != nil {
		return err
	}
	if exists && !cfg.reset {
		return fmt.Errorf("seed data already exists; pass --reset to replace data for %s", strings.Join(hostnames, ", "))
	}
	if cfg.reset {
		if err := resetSeedData(ctx, appDB, analyticsDB, hostnames); err != nil {
			return err
		}
	}

	now := time.Now().UTC().Truncate(time.Hour)
	if err := createWebsites(ctx, appDB, user.ID, hostnames, now); err != nil {
		return err
	}

	viewCounts := splitViews(views)
	totalEvents, err := seedAnalytics(ctx, analyticsDB, viewCounts, cfg.days, now)
	if err != nil {
		return err
	}

	fmt.Printf("Seeded %s dev data into %s and %s\n", cfg.size, cfg.appDB, cfg.analyticsDB)
	fmt.Printf("User: %s\n", user.Username)
	for i, site := range sites {
		fmt.Printf("- %s: %d views\n", site.hostname, viewCounts[i])
	}
	fmt.Printf("Total: %d views, %d custom properties\n", views, totalEvents)
	return nil
}

func getSeedOwner(ctx context.Context, appDB *sqlite.Client) (*seedOwner, error) {
	var user seedOwner
	err := appDB.QueryRowxContext(ctx, `--sql
		SELECT id, username
		FROM users
		ORDER BY date_created, id
		LIMIT 1`,
	).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("get seed owner: %w", err)
	}
	return &user, nil
}

func hasSeedData(
	ctx context.Context,
	appDB *sqlite.Client,
	analyticsDB *duckdb.Client,
	hostnames []string,
) (bool, error) {
	for _, hostname := range hostnames {
		if exists, err := rowExists(ctx, appDB.DB, "SELECT 1 FROM websites WHERE hostname = ? LIMIT 1", hostname); err != nil {
			return false, fmt.Errorf("check website %q: %w", hostname, err)
		} else if exists {
			return true, nil
		}

		if exists, err := rowExists(ctx, analyticsDB.DB, "SELECT 1 FROM views WHERE hostname = ? LIMIT 1", hostname); err != nil {
			return false, fmt.Errorf("check views %q: %w", hostname, err)
		} else if exists {
			return true, nil
		}

		if exists, err := rowExists(ctx, analyticsDB.DB, "SELECT 1 FROM events WHERE group_name = ? LIMIT 1", hostname); err != nil {
			return false, fmt.Errorf("check events %q: %w", hostname, err)
		} else if exists {
			return true, nil
		}
	}
	return false, nil
}

func rowExists(ctx context.Context, db *sqlx.DB, query string, args ...any) (bool, error) {
	var marker int
	err := db.QueryRowxContext(ctx, query, args...).Scan(&marker)
	switch {
	case err == nil:
		return true, nil
	case err == sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func resetSeedData(
	ctx context.Context,
	appDB *sqlite.Client,
	analyticsDB *duckdb.Client,
	hostnames []string,
) error {
	for _, hostname := range hostnames {
		if _, err := analyticsDB.ExecContext(ctx, "DELETE FROM events WHERE group_name = ?", hostname); err != nil {
			return fmt.Errorf("delete events for %q: %w", hostname, err)
		}
		if _, err := analyticsDB.ExecContext(ctx, "DELETE FROM views WHERE hostname = ?", hostname); err != nil {
			return fmt.Errorf("delete views for %q: %w", hostname, err)
		}
		if _, err := appDB.ExecContext(ctx, "DELETE FROM websites WHERE hostname = ?", hostname); err != nil {
			return fmt.Errorf("delete website %q: %w", hostname, err)
		}
	}
	return nil
}

func createWebsites(
	ctx context.Context,
	appDB *sqlite.Client,
	userID string,
	hostnames []string,
	now time.Time,
) error {
	for _, hostname := range hostnames {
		website := model.NewWebsite(userID, hostname, now.Unix(), now.Unix())
		if err := appDB.CreateWebsite(ctx, website); err != nil {
			return fmt.Errorf("create website %q: %w", hostname, err)
		}
	}
	return nil
}

func splitViews(total int) []int {
	counts := make([]int, len(sites))
	weightTotal := 0
	for _, site := range sites {
		weightTotal += site.weight
	}

	allocated := 0
	for i, site := range sites {
		counts[i] = total * site.weight / weightTotal
		allocated += counts[i]
	}
	counts[0] += total - allocated
	return counts
}

func seedAnalytics(
	ctx context.Context,
	analyticsDB *duckdb.Client,
	viewCounts []int,
	days int,
	end time.Time,
) (int, error) {
	tx, err := analyticsDB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin analytics transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Commit closes the transaction on success.

	// The repository methods stamp current time; the dashboard seed needs recent historical rows.
	viewStmt, err := tx.PreparexContext(ctx, `--sql
		INSERT INTO views (
			bid,
			hostname,
			pathname,
			is_unique_user,
			is_unique_page,
			referrer_host,
			referrer_group,
			country,
			language_base,
			language_dialect,
			ua_browser,
			ua_os,
			ua_device_type,
			utm_source,
			utm_medium,
			utm_campaign,
			duration_ms,
			date_created
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("prepare views insert: %w", err)
	}
	defer viewStmt.Close()

	eventStmt, err := tx.PreparexContext(ctx, `--sql
		INSERT INTO events (
			bid,
			batch_id,
			group_name,
			name,
			value,
			date_created
		) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("prepare events insert: %w", err)
	}
	defer eventStmt.Close()

	totalEvents := 0
	for siteIndex, site := range sites {
		for i := range viewCounts[siteIndex] {
			view, durationMS, dateCreated := generatePageView(site, siteIndex, i, days, end)
			if _, err := viewStmt.ExecContext(ctx,
				view.BID,
				view.Hostname,
				view.Pathname,
				view.IsUniqueUser,
				view.IsUniquePage,
				view.ReferrerHost,
				view.ReferrerGroup,
				view.Country,
				view.LanguageBase,
				view.LanguageDialect,
				view.BrowserName,
				view.OS,
				view.DeviceType,
				view.UTMSource,
				view.UTMMedium,
				view.UTMCampaign,
				durationMS,
				dateCreated,
			); err != nil {
				return 0, fmt.Errorf("insert view %q: %w", view.BID, err)
			}

			events := generateEvents(view, i)
			for _, event := range events {
				if _, err := eventStmt.ExecContext(ctx,
					event.BID,
					event.BatchID,
					event.Group,
					event.Name,
					event.Value,
					dateCreated.Add(time.Duration(i%180)*time.Second),
				); err != nil {
					return 0, fmt.Errorf("insert event %q: %w", event.BatchID, err)
				}
			}
			totalEvents += len(events)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit analytics transaction: %w", err)
	}
	return totalEvents, nil
}

func generatePageView(site site, siteIndex int, index int, days int, end time.Time) (*model.PageViewHit, int, time.Time) {
	referrer := pickString(referrers, index+siteIndex)
	language := languages[index%len(languages)]
	durationMS := 3_000 + (index*7_919)%240_000
	dateCreated := end.Add(-time.Duration((index*37+siteIndex*263)%(days*24*60)) * time.Minute)

	view := &model.PageViewHit{
		BID:             fmt.Sprintf("seed-%d-%d", siteIndex, index),
		Hostname:        site.hostname,
		Pathname:        pickString(site.paths, index),
		IsUniqueUser:    index%3 == 0,
		IsUniquePage:    index%2 == 0,
		ReferrerHost:    referrer,
		ReferrerGroup:   referrer,
		Country:         pickString(countries, index+siteIndex),
		LanguageBase:    language.base,
		LanguageDialect: language.dialect,
		BrowserName:     pickString(browsers, index),
		OS:              pickString(operatingSystems, index/2),
		DeviceType:      pickString(devices, index/3),
	}

	if index%3 == 0 {
		view.UTMSource = pickString(utmSources, index)
		view.UTMMedium = pickString(utmMediums, index)
		view.UTMCampaign = pickString(utmCampaigns, index)
	}

	return view, durationMS, dateCreated
}

func generateEvents(view *model.PageViewHit, index int) []model.EventHit {
	if index%4 != 0 {
		return nil
	}

	batchID := fmt.Sprintf("event-%s", view.BID)
	events := make([]model.EventHit, 0, len(properties))
	for i, property := range properties {
		if i > 0 && index%(i+2) != 0 {
			continue
		}
		events = append(events, model.EventHit{
			BID:     view.BID,
			BatchID: batchID,
			Group:   view.Hostname,
			Name:    property.name,
			Value:   pickString(property.values, index+i),
		})
	}
	return events
}

func pickString(items []string, index int) string {
	return items[index%len(items)]
}
